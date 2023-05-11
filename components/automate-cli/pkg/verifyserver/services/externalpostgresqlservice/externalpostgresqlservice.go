package externalpostgresqlservice

import (
	"fmt"

	//"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/utils/db"
	"github.com/chef/automate/lib/logger"
	"github.com/pkg/errors"
)

var (
	ExternalPgSuccessConnectionTitle           = "Connection successfully tested"
	ExternalPgFailConnectionTitle           = "External Postgresql Connection failed"
	ExternalPgConnectionErrorMsg        = "Machine is unable to connect with External Managed Postgresql"
	ExternalPgConnectionResolutionMsg   = "Ensure that the Postgres configuration provided is correct. Review security group or firewall settings as well on the infrastructure"
	ExternalPgConnectionSuccessMsg      = "Connection successfully tested"
)
type ISExternalPostgresqlService interface {
	GetPgConnection(models.ExternalPgRequest) models.ExternalPgConnectionDetails
}

type ExternalPostgresqlService struct{
	logger logger.Logger
	DBUtils db.DB
	req models.ExternalPgRequest
}

func NewExternalPostgresqlService(db db.DB,logger logger.Logger) ISExternalPostgresqlService {
	return &ExternalPostgresqlService{
		logger: logger,
		DBUtils:db,
	}
}

func (pg *ExternalPostgresqlService) GetPgConnection(req models.ExternalPgRequest) models.ExternalPgConnectionDetails {
	pg.req = req
	err := pg.CheckExternalPgConnection()
	if err != nil {
		return pg.Response(ExternalPgFailConnectionTitle, "", errors.Wrap(err, ExternalPgConnectionErrorMsg).Error(), ExternalPgConnectionResolutionMsg, false)
		// 	Title:         "Postgres Connection failed",
		// 	Passed:        false,
		// 	Status:		   "PASS",
		// 	SuccessMsg:    "",
		// 	ErrorMsg:      "Machine is unable to connect with External Managed Postgresql",
		// 	ResolutionMsg: "Ensure that the Postgres configuration provided is correct. Review security group or firewall settings as well on the infrastructure",			
		// }
	}

	return pg.Response(ExternalPgSuccessConnectionTitle,ExternalPgConnectionSuccessMsg,"","", true)
	// 	Title:         "Connection successfully tested",
	// 	Passed:        true,
	// 	Status:		   "PASS",
	// 	SuccessMsg:    "Machine is able to connect with External Managed Postgres",
	// 	ErrorMsg:      "",
	// 	ResolutionMsg: "",
	// }
}


	
func (p *ExternalPostgresqlService) CheckExternalPgConnection() error {
		//construct the connection string
		con := fmt.Sprintf("PostgresqlInstanceUrl=%s PostgresqlSuperUserUserName=%s PostgresqlSuperUserPassword=%s PostgresqlDbUserUserName=%s PostgresqlDbUserPassword=%s sslrootcert=%s",p.req.PostgresqlInstanceUrl,p.req.PostgresqlSuperUserUserName,p.req.PostgresqlSuperUserPassword,p.req.PostgresqlDbUserUserName,p.req.PostgresqlDbUserPassword,p.req.PostgresqlRootCert)
		err := p.DBUtils.InitPostgresDB(con)
		if err != nil{
			p.logger.Error("External Postgresql Connection failed: ", err.Error())
			return err
		}
		p.logger.Info("External Postgresql aws connection success")
		return nil
}

func (pg *ExternalPostgresqlService) Response(Title, SuccessMsg, ErrorMsg, ResolutionMsg string, Passed bool) models.ExternalPgConnectionDetails {
	return models.ExternalPgConnectionDetails{
		Title:         Title,
		Passed:        Passed,
		SuccessMsg:    SuccessMsg,
		ErrorMsg:      ErrorMsg,
		ResolutionMsg: ResolutionMsg,
	}
}