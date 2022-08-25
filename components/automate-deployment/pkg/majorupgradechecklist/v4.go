package majorupgradechecklist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/chef/automate/components/automate-cli/pkg/status"
	"github.com/chef/automate/components/automate-deployment/pkg/cli"
	"github.com/pkg/errors"
)

type indexVersion struct {
	Settings struct {
		Index struct {
			Version struct {
				CreatedString string `json:"created_string"`
				Created       string `json:"created"`
			} `json:"version"`
		} `json:"index"`
	} `json:"settings"`
}

type indexDetails struct {
	Name    string
	Version string
}

type V4ChecklistManager struct {
	writer       cli.FormatWriter
	version      string
	isExternalES bool
}

const (
	initMsgV4 = `This is a Major upgrade. 
========================

  1) In this release Elastic Search is migrated to Open Search
  2) This upgrade will require special care. 

===== Your installation is using %s Elastic Search =====

  Please confirm this checklist that you have taken care of these steps 
  before continuing with the Upgrade to version %s:
`

	externalESUpgradeMsg = "You are ready with your external open search with migrated data from elastic search, this should be done with the help of your Administrator"

	externalESUpgradeError = "As you are using external elastic search, migrations of data to open search is required before the upgrade"

	shardingError = "sharding has to be disabled before upgrade"

	run_os_data_migrate = `Migrate Data from Elastic Search to Open Search using this command:
     $ ` + run_os_data_migrate_cmd

	run_os_data_migrate_cmd = `chef-automate post-major-upgrade migrate --data=es`

	run_os_data_cleanup = `If you are sure all data is available in Upgraded Automate, then we can free up old elastic search Data by running: 
     $ ` + run_os_data_cleanup_cmd

	run_os_data_cleanup_cmd = `chef-automate post-major-upgrade clear-data --data=es`

	disable_maintenance_mode = `Disable the maintenance mode if you enabled previously using:
	$ ` + disable_maintenance_mode_cmd

	disable_maintenance_mode_cmd = `chef-automate maintenance off`
)

var postChecklistV4Embedded = []PostCheckListItem{
	{
		Id:         "upgrade_status",
		Msg:        run_chef_automate_upgrade_status,
		Cmd:        run_chef_automate_upgrade_status_cmd,
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "disable_maintenance_mode",
		Msg:        disable_maintenance_mode,
		Cmd:        disable_maintenance_mode_cmd,
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "migrate_es",
		Msg:        run_os_data_migrate,
		Cmd:        run_os_data_migrate_cmd,
		IsExecuted: false,
	}, {
		Id:         "check_ui",
		Msg:        ui_check,
		Cmd:        "",
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "clean_up",
		Msg:        run_os_data_cleanup,
		Cmd:        run_os_data_cleanup_cmd,
		Optional:   true,
		IsExecuted: false,
	},
}

var postChecklistV4External = []PostCheckListItem{
	{
		Id:         "upgrade_status",
		Msg:        run_chef_automate_upgrade_status,
		Cmd:        run_chef_automate_upgrade_status_cmd,
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "status",
		Msg:        run_chef_automate_status,
		Cmd:        run_chef_automate_status_cmd,
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "disable_maintenance_mode",
		Msg:        disable_maintenance_mode,
		Cmd:        disable_maintenance_mode_cmd,
		Optional:   true,
		IsExecuted: false,
	}, {
		Id:         "check_ui",
		Msg:        ui_check,
		Cmd:        "",
		Optional:   true,
		IsExecuted: false,
	},
}

func NewV4ChecklistManager(writer cli.FormatWriter, version string) *V4ChecklistManager {
	return &V4ChecklistManager{
		writer:       writer,
		version:      version,
		isExternalES: IsExternalElasticSearch(writer),
	}
}

func (ci *V4ChecklistManager) GetPostChecklist() []PostCheckListItem {
	var postChecklist []PostCheckListItem
	if ci.isExternalES {
		postChecklist = postChecklistV4External
	} else {
		postChecklist = postChecklistV4Embedded
	}
	return postChecklist
}

func (ci *V4ChecklistManager) RunChecklist() error {
	var dbType string
	checklists := []Checklist{}
	var postcheck []PostCheckListItem

	if ci.isExternalES {
		dbType = "External"
		postcheck = postChecklistV4External
		checklists = append(checklists, []Checklist{downTimeCheckV4(), backupCheck(), externalESUpgradeCheck(),
			postChecklistIntimationCheckV4(!ci.isExternalES)}...)
	} else {
		dbType = "Embedded"
		postcheck = postChecklistV4Embedded
		checklists = append(checklists, []Checklist{runIndexCheck(), downTimeCheckV4(), backupCheck(), diskSpaceCheck(),
			disableSharding(), postChecklistIntimationCheckV4(!ci.isExternalES)}...)
	}
	checklists = append(checklists, showPostChecklist(&postcheck), promptUpgradeContinueV4(!ci.isExternalES))

	helper := ChecklistHelper{
		Writer: ci.writer,
	}

	ci.writer.Println(fmt.Sprintf(initMsgV4, dbType, ci.version)) //display the init message

	for _, item := range checklists {
		if item.TestFunc == nil {
			continue
		}
		if err := item.TestFunc(helper); err != nil {
			return errors.Wrap(err, "one of the checklist was not accepted/satisfied for upgrade")
		}
	}
	return nil
}

func postChecklistIntimationCheckV4(isEmbedded bool) Checklist {
	return Checklist{
		Name:        "post_checklist_intimation_acceptance",
		Description: "confirmation check for post checklist intimation",
		TestFunc: func(h ChecklistHelper) error {
			resp, err := h.Writer.Confirm("After this upgrade completes, you will have to run Post upgrade steps to ensure your data is migrated and your Automate is ready for use")
			if err != nil {
				h.Writer.Error(err.Error())

				shardError := enableSharding(h, isEmbedded)
				if shardError != nil {
					h.Writer.Error(shardError.Error())
				}

				return status.Errorf(status.InvalidCommandArgsError, err.Error())
			}
			if !resp {
				h.Writer.Error(postChecklistIntimationError)

				shardError := enableSharding(h, isEmbedded)
				if shardError != nil {
					h.Writer.Error(shardError.Error())
				}

				return status.New(status.InvalidCommandArgsError, postChecklistIntimationError)
			}
			return nil
		},
	}
}

func promptUpgradeContinueV4(isEmbedded bool) Checklist {
	return Checklist{
		Name:        "post_checklist",
		Description: "display post checklist and ask for final confirmation",
		TestFunc: func(h ChecklistHelper) error {
			resp, err := h.Writer.Confirm(v3_post_checklist_confirmation)
			if err != nil {
				h.Writer.Error(err.Error())

				shardError := enableSharding(h, isEmbedded)
				if shardError != nil {
					h.Writer.Error(shardError.Error())
				}

				return status.Errorf(status.InvalidCommandArgsError, err.Error())
			}
			if !resp {
				h.Writer.Error("end user not ready to upgrade")
				shardError := enableSharding(h, isEmbedded)
				if shardError != nil {
					h.Writer.Error(shardError.Error())
				}
				return status.New(status.InvalidCommandArgsError, "end user not ready to upgrade")
			}
			return nil
		},
	}
}

func enableSharding(h ChecklistHelper, isEmbedded bool) error {
	if isEmbedded {
		h.Writer.Println("enabling sharding")
		return reEnableShardAllocation()
	}
	return nil
}

func PatchBestOpenSearchSettings(isEmbedded bool) error {
	if isEmbedded {
		bestSettings, err := GetBestSettings()
		if err != nil {
			return err
		}
		fmt.Println(bestSettings)
		temp := template.Must(template.New("init").
			Funcs(template.FuncMap{"StringsJoin": strings.Join}).
			Parse(patchOpensearchSettingsTemplate))

		var buf bytes.Buffer
		err = temp.Execute(&buf, bestSettings)
		if err != nil {
			fmt.Printf("some error occurred while rendering template \n %s", err)
		}
		finalTemplate := buf.String()
		err = ioutil.WriteFile(AutomateOpensearchConfigPatch, []byte(finalTemplate), 0600)
		if err != nil {
			fmt.Printf("Writing into %s file failed\n", AutomateOpensearchConfigPatch)
		}
		fmt.Printf("Config written to %s\n", AutomateOpensearchConfigPatch)
		cmd := exec.Command("chef-automate", "config", "patch", AutomateOpensearchConfigPatch)
		cmd.Stdin = os.Stdin
		cmd.Stdout = io.MultiWriter(os.Stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error in patching opensearch configuration")
			fmt.Println(err)
		}
		fmt.Printf("config patch execution done, exiting\n")
	}
	return nil
}

func GetBestSettings() (*ESSettings, error) {
	esSetting, err := getOldElasticSearchSettings()
	if err != nil {
		fmt.Println(err)
	}
	ossSetting := GetDefaultOSSettings()
	bestSetting := &ESSettings{}

	esHeapSize, err := extractNumericFromText(esSetting.HeapMemory, 0)
	if err != nil {
		fmt.Println("Not able to parse heapsize from elasticsearch settings")
	} else {
		ossHeapSize, err := extractNumericFromText(ossSetting.HeapMemory, 0)
		if err != nil {
			fmt.Println("Not able to parse heapsize from opensearch settings")
		} else {
			if esHeapSize >= 32 || ossHeapSize >= 32 {
				bestSetting.HeapMemory = "32g"
			}
			if esHeapSize > ossHeapSize {
				bestSetting.HeapMemory = fmt.Sprintf("%d", int64(math.Round(esHeapSize)))
			} else {
				bestSetting.HeapMemory = fmt.Sprintf("%d", int64(math.Round(ossHeapSize)))
			}
		}
	}

	esIndicesBreakerTotalLimit, err := extractNumericFromText(esSetting.IndicesBreakerTotalLimit, 0)
	if err != nil {
		fmt.Println("Not able to parse indicesBreakerTotalLimit from opensearch settings")
	} else {
		ossIndicesBreakerTotalLimit, err := extractNumericFromText(ossSetting.IndicesBreakerTotalLimit, 0)
		if err != nil {
			fmt.Println("Not able to parse indicesBreakerTotalLimit from opensearch settings")
		} else {
			if esIndicesBreakerTotalLimit > ossIndicesBreakerTotalLimit {
				bestSetting.IndicesBreakerTotalLimit = fmt.Sprintf("%d", int64(math.Round(esIndicesBreakerTotalLimit)))
			} else {
				bestSetting.IndicesBreakerTotalLimit = fmt.Sprintf("%d", int64(math.Round(ossIndicesBreakerTotalLimit)))
			}
		}
	}

	esRunTimeMaxOpenFile, err := extractNumericFromText(esSetting.RuntimeMaxOpenFile, 0)
	if err != nil {
		fmt.Println("Not able to parse esRunTimeMaxOpenFile from opensearch settings")
	} else {
		ossRunTimeMaxOpenFile, err := extractNumericFromText(ossSetting.RuntimeMaxOpenFile, 0)
		if err != nil {
			fmt.Println("Not able to parse ossRunTimeMaxOpenFile from opensearch settings")
		} else {
			if esRunTimeMaxOpenFile > ossRunTimeMaxOpenFile {
				bestSetting.RuntimeMaxOpenFile = fmt.Sprintf("%d", int64(esRunTimeMaxOpenFile))
			} else {
				bestSetting.RuntimeMaxOpenFile = fmt.Sprintf("%d", int64(ossRunTimeMaxOpenFile))
			}
		}
	}

	if esSetting.TotalShardSettings > ossSetting.TotalShardSettings {
		bestSetting.TotalShardSettings = esSetting.TotalShardSettings
	} else {
		bestSetting.TotalShardSettings = ossSetting.TotalShardSettings
	}

	if esSetting.RuntimeMaxLockedMem == "unlimited" || ossSetting.RuntimeMaxLockedMem == "unlimited" {
		bestSetting.RuntimeMaxLockedMem = "unlimited"
	} else {
		esRunTimeMaxLockedMem, err := extractNumericFromText(esSetting.RuntimeMaxLockedMem, 0)
		if err != nil {
			fmt.Println("Not able to parse esRunTimeMaxLockedMem from opensearch settings")
		} else {
			ossRunTimeMaxLockedMem, err := extractNumericFromText(ossSetting.RuntimeMaxLockedMem, 0)
			if err != nil {
				fmt.Println("Not able to parse ossRunTimeMaxLockedMem from opensearch settings")
			} else {
				if esRunTimeMaxLockedMem > ossRunTimeMaxLockedMem {
					bestSetting.RuntimeMaxLockedMem = fmt.Sprintf("%d", int64(math.Round(esRunTimeMaxLockedMem)))
				} else {
					bestSetting.RuntimeMaxLockedMem = fmt.Sprintf("%d", int64(math.Round(ossRunTimeMaxLockedMem)))
				}
			}
		}
	}
	return bestSetting, nil
}

func extractNumericFromText(text string, index int) (float64, error) {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	submatchall := re.FindAllString(text, -1)
	if len(submatchall) < 1 {
		return 0, errors.New("No match found")
	}
	//fmt.Println(submatchall)
	numeric, err := strconv.ParseFloat(submatchall[index], 64)
	if err != nil {
		return 0, err
	}
	return numeric, nil
}

func getOldElasticSearchSettings() (*ESSettings, error) {
	esSetting := &ESSettings{}
	jsonData, err := ioutil.ReadFile(V3ESSettingFile)
	if err != nil {
		return esSetting, err
	}
	err = json.Unmarshal(jsonData, esSetting)
	if err != nil {
		return esSetting, err
	}
	return esSetting, nil
}

func GetDefaultOSSettings() *ESSettings {
	defaultSettings := &ESSettings{}
	defaultSettings.HeapMemory = defaultHeapSize()
	defaultSettings.IndicesBreakerTotalLimit = INDICES_BREAKER_TOTAL_LIMIT_DEFAULT
	defaultSettings.RuntimeMaxLockedMem = MAX_LOCKED_MEM_DEFAULT
	defaultSettings.RuntimeMaxOpenFile = MAX_OPEN_FILE_DEFAULT
	defaultSettings.TotalShardSettings = INDICES_TOTAL_SHARD_DEFAULT
	return defaultSettings
}

func downTimeCheckV4() Checklist {
	return Checklist{
		Name:        "down_time_acceptance",
		Description: "confirmation for downtime",
		TestFunc: func(h ChecklistHelper) error {
			resp, err := h.Writer.Confirm("You had planned for a downtime, by running the command(chef-automate maintenance on)?:")
			if err != nil {
				h.Writer.Error(err.Error())
				return status.Errorf(status.InvalidCommandArgsError, err.Error())
			}
			if !resp {
				h.Writer.Error(downTimeError)
				return status.New(status.InvalidCommandArgsError, downTimeError)
			}
			return nil
		},
	}
}

func externalESUpgradeCheck() Checklist {
	return Checklist{
		Name:        "external_es_upgrade_acceptance",
		Description: "confirmation check for external es upgrade",
		TestFunc: func(h ChecklistHelper) error {
			resp, err := h.Writer.Confirm(externalESUpgradeMsg)
			if err != nil {
				h.Writer.Error(err.Error())
				return status.Errorf(status.InvalidCommandArgsError, err.Error())
			}
			if !resp {
				h.Writer.Error(externalESUpgradeError)
				return status.New(status.InvalidCommandArgsError, externalESUpgradeError)
			}
			return nil
		},
	}
}

/*
const disableShardCurlRequest = `curl --location --request PUT 'http://localhost:10141/_cluster/settings'
-u 'admin:admin' --header 'Content-Type: application/json'
--data-raw '{ "persistent": { "cluster.routing.allocation.enable": "primaries"}}'`

const flushIndexRequest = `curl --location --request PUT 'http://localhost:10168/_cluster/settings' \
-u 'admin:admin' \
--header 'Content-Type: application/json' \
--data-raw '{
  "persistent": {
    "cluster.routing.allocation.enable": "all"
  }}'`

*/
/*
disableShardAllocation : Disable shard allocation: during upgrade
we do not want any shard movement as the cluster will be
taken down.
*/
func disableShardAllocation() error {
	url := "http://localhost:10141/_cluster/settings"
	method := "PUT"

	payload := strings.NewReader(`{
    "persistent": {
        "cluster.routing.allocation.enable": "primaries"
    }
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload) // nosemgrep

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body) // nosemgrep
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

/*
flushRequest : Stop indexing, and perform a flush:
as the cluster will be taken down, indexing/searching
should be stopped and _flush can be used to permanently
store information into the index which will prevent any
data loss during upgrade
*/
func flushRequest() error {

	url := "http://localhost:10141/_flush"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil) // nosemgrep

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body) // nosemgrep
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

/*
reEnableShardAllocation: Re-enable shard allocation
In case of Any error we can Undo the First request
*/
func reEnableShardAllocation() error {
	url := "http://localhost:10141/_cluster/settings"
	method := "PUT"

	payload := strings.NewReader(`{
		"persistent": {
			"cluster.routing.allocation.enable": null
		  }
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload) // nosemgrep

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body) // nosemgrep
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func disableSharding() Checklist {
	return Checklist{
		Name:        "disable sharding",
		Description: "confirmation check to disable sharding",
		TestFunc: func(h ChecklistHelper) error {
			resp, err := h.Writer.Confirm("This will disable Sharding on your elastic search")
			if err != nil {
				h.Writer.Error(err.Error())
				return status.Errorf(status.InvalidCommandArgsError, err.Error())
			}
			if !resp {
				h.Writer.Error(shardingError)
				return status.New(status.InvalidCommandArgsError, shardingError)
			}
			err = disableShardAllocation()
			if err != nil {
				h.Writer.Error(err.Error())
				return status.Errorf(status.DatabaseError, err.Error())
			}

			err = flushRequest()
			if err != nil {
				h.Writer.Error(err.Error())
				return status.Errorf(status.DatabaseError, err.Error())
			}
			return nil
		},
	}
}

const habrootcmd = "HAB_LICENSE=accept-no-persist hab pkg path chef/deployment-service"

func getHabRootPath(habrootcmd string) string {
	out, err := exec.Command("/bin/sh", "-c", habrootcmd).Output()
	if err != nil {
		return "/hab/"
	}
	pkgPath := string(out) // /a/b/c/hab    /hab/svc
	habIndex := strings.Index(string(pkgPath), "hab")
	rootHab := pkgPath[0 : habIndex+4] // this will give <>/<>/hab/
	if rootHab == "" {
		rootHab = "/hab/"
	}
	return rootHab
}

func checkIndexVersion() error {
	var basePath = "http://localhost:10144/"

	habpath := getHabRootPath(habrootcmd)

	input, err := ioutil.ReadFile(habpath + "svc/automate-es-gateway/config/URL") // nosemgrep
	if err != nil {
		fmt.Printf("Failed to read URL file")
	}
	url := strings.TrimSuffix(string(input), "\n")
	if url != "" {
		basePath = "http://" + url + "/"
	}

	allIndexList, err := getDataFromUrl(basePath + "_cat/indices?h=index")
	if err != nil {
		return err
	}

	indexDetailsArray := []indexDetails{}
	for _, index := range strings.Split(strings.TrimSuffix(string(allIndexList), "\n"), "\n") {
		versionData, err := getDataFromUrl(basePath + index + "/_settings/index.version.created*?&human")
		if err != nil {
			return err
		}
		i, createdString, err := getMajorVersion(versionData, index)
		if err != nil {
			return err
		}
		if i < 6 {
			indexDetailsArray = append(indexDetailsArray, indexDetails{Name: index, Version: createdString})
		}
	}
	if len(indexDetailsArray) > 0 {
		return formErrorMsg(indexDetailsArray)
	}
	return nil
}

func formErrorMsg(IndexDetailsArray []indexDetails) error {
	msg := "\nUnsupported index versions. To continue with the upgrade, please reindex the indices shown below to version 6.\n"
	for _, version := range IndexDetailsArray {
		msg += fmt.Sprintf("- Index Name: %s, Version: %s \n", version.Name, version.Version)
	}
	msg += "\nFollow the guide below to learn more about reindexing:\nhttps://www.elastic.co/guide/en/elasticsearch/reference/6.8/docs-reindex.html"
	return fmt.Errorf(msg)
}

func getMajorVersion(versionData []byte, index string) (int64, string, error) {
	data := map[string]interface{}{}
	json.Unmarshal(versionData, &data)
	dataIdx := indexVersion{}
	b, err := json.Marshal(data[index])
	if err != nil {
		return -1, "", errors.Wrap(err, "failed to marshal index data")
	}
	err = json.Unmarshal(b, &dataIdx)
	if err != nil {
		return -1, "", errors.Wrap(err, "failed to unmarshal index data")
	}
	if dataIdx.Settings.Index.Version.CreatedString == "" {
		return -1, "", errors.New("version not found for index")
	}
	i, err := strconv.ParseInt(dataIdx.Settings.Index.Version.CreatedString[0:1], 10, 64)
	if err != nil {
		return -1, "", errors.Wrap(err, "failed to parse index version")
	}
	return i, dataIdx.Settings.Index.Version.CreatedString, nil
}

func runIndexCheck() Checklist {
	return Checklist{
		Name:        "check index version",
		Description: "confirmation check index version",
		TestFunc: func(h ChecklistHelper) error {
			err := checkIndexVersion()
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func StoreESSettings() error {
	pid, err := getElasticsearchPID()
	if err != nil {
		fmt.Println("No process id found for running elasticsearch")
	}
	esSettings, err := getAllSearchEngineSettings(pid)
	if err != nil {
		fmt.Println(err)
	}
	esSettingsJson, err := json.Marshal(esSettings)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(esSettingsJson))
	err = ioutil.WriteFile(V3ESSettingFile, esSettingsJson, 775)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (ci *V4ChecklistManager) StoreSearchEngineSettings() error {
	return StoreESSettings()
}
