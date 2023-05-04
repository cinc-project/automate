package externalpostgresqlservice

import (
	"fmt"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/utils/db"
)

type ISExternalPostgresqlService interface {
	GetPgConnection(models.ExternalPgRequest) models.ExternalPgConnectionDetails
}

type ExternalPostgresqlService struct{
	DBUtils db.DB
	req *models.ExternalPgRequest
}

func NewExternalPostgresqlService(db db.DB) ISExternalPostgresqlService {
	return &ExternalPostgresqlService{DBUtils:db}
}

func (pg *ExternalPostgresqlService) GetPgConnection(req models.ExternalPgRequest) models.ExternalPgConnectionDetails {
	pg.req = &req
	err := pg.CheckExternalPgConnection()
	if err != nil {
		return models.ExternalPgConnectionDetails{
			Title:         "Postgres Connection failed",
			Passed:        false,
			Status:		   "PASS",
			SuccessMsg:    "",
			ErrorMsg:      "Machine is unable to connect with External Managed Postgresql",
			ResolutionMsg: "Ensure that the Postgres configuration provided is correct. Review security group or firewall settings as well on the infrastructure",			
		}
	}

	return models.ExternalPgConnectionDetails{
		Title:         "Connection successfully tested",
		Passed:        true,
		Status:		   "PASS",
		SuccessMsg:    "Machine is able to connect with External Managed Postgres",
		ErrorMsg:      "",
		ResolutionMsg: "",
	}
}


	
func (p *ExternalPostgresqlService) CheckExternalPgConnection() error {
		//construct the connection string
		con := fmt.Sprintf("PostgresqlInstanceUrl=%s PostgresqlSuperUserUserName=%s PostgresqlSuperUserPassword=%s PostgresqlDbUserUserName=%s PostgresqlDbUserPassword=%s sslrootcert=%s",p.req.PostgresqlInstanceUrl,p.req.PostgresqlSuperUserUserName,p.req.PostgresqlSuperUserPassword,p.req.PostgresqlDbUserUserName,p.req.PostgresqlDbUserPassword,p.req.PostgresqlRootCert)
		err := p.DBUtils.InitPostgresDB(con)
		if err != nil{
			return err
		}
		return nil
}
