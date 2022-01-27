package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-process-order-conf-reads-rmq-kube/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"
	"golang.org/x/xerrors"
)

type RMQOutputter interface {
	Send(sendQueue string, payload map[string]interface{}) error
}

type SAPAPICaller struct {
	baseURL      string
	apiKey       string
	outputQueues []string
	outputter    RMQOutputter
	log          *logger.Logger
}

func NewSAPAPICaller(baseUrl string, outputQueueTo []string, outputter RMQOutputter, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:      baseUrl,
		apiKey:       GetApiKey(),
		outputQueues: outputQueueTo,
		outputter:    outputter,
		log:          l,
	}
}

func (c *SAPAPICaller) AsyncGetProcessOrderConfirmation(orderID, batch, confirmationGroup string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "ConfByOrderID":
			func() {
				c.ConfByOrderID(orderID)
				wg.Done()
			}()
		case "MaterialMovements":
			func() {
				c.MaterialMovements(batch)
				wg.Done()
			}()
		case "BatchCharacteristic":
			func() {
				c.BatchCharacteristic(batch)
				wg.Done()
			}()
		case "ConfByOrderIDConfGroup":
			func() {
				c.ConfByOrderIDConfGroup(orderID, confirmationGroup)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) ConfByOrderID(orderID string) {
	confByOrderIDData, err := c.callProcessOrderConfirmationSrvAPIRequirementConfByOrderID("ProcOrdConf2", orderID)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": confByOrderIDData, "function": "ProcessOrderConfirmationConfirmation"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(confByOrderIDData)

}

func (c *SAPAPICaller) callProcessOrderConfirmationSrvAPIRequirementConfByOrderID(api, orderID string) ([]sap_api_output_formatter.Confirmation, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithConfByOrderID(req, orderID)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToConfirmation(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) MaterialMovements(batch string) {
	materialMovementsData, err := c.callProcessOrderConfirmationSrvAPIRequirementMaterialMovements("ProcOrdConfMatlDocItm", batch)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": materialMovementsData, "function": "ProcessOrderConfirmationMaterialMovements"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(materialMovementsData)
}

func (c *SAPAPICaller) callProcessOrderConfirmationSrvAPIRequirementMaterialMovements(api, batch string) ([]sap_api_output_formatter.MaterialMovements, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithMaterialMovements(req, batch)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToMaterialMovements(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) ConfByOrderIDConfGroup(orderID, confirmationGroup string) {
	confByOrderIDConfGroupData, err := c.callProcessOrderConfirmationSrvAPIRequirementConfByOrderIDConfGroup("ProcOrdConf2", orderID, confirmationGroup)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": confByOrderIDConfGroupData, "function": "ProcessOrderConfirmationConfirmation"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(confByOrderIDConfGroupData)

}

func (c *SAPAPICaller) callProcessOrderConfirmationSrvAPIRequirementConfByOrderIDConfGroup(api, orderID, confirmationGroup string) ([]sap_api_output_formatter.Confirmation, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithConfByOrderIDConfGroup(req, orderID, confirmationGroup)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToConfirmation(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) BatchCharacteristic(batch string) {
	batchCharacteristicData, err := c.callProcessOrderConfirmationSrvAPIRequirementBatchCharacteristic("ProcOrderConfBatchCharc", batch)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": batchCharacteristicData, "function": "ProcessOrderConfirmationBatchCharacteristic"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(batchCharacteristicData)

}

func (c *SAPAPICaller) callProcessOrderConfirmationSrvAPIRequirementBatchCharacteristic(api, batch string) ([]sap_api_output_formatter.BatchCharacteristic, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithBatchCharacteristic(req, batch)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToBatchCharacteristic(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithConfByOrderID(req *http.Request, orderID string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("OrderID eq '%s'", orderID))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithMaterialMovements(req *http.Request, batch string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Batch eq '%s'", batch))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithConfByOrderIDConfGroup(req *http.Request, orderID, confirmationGroup string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("OrderID eq '%s' and ConfirmationGroup eq '%s'", orderID, confirmationGroup))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithBatchCharacteristic(req *http.Request, batch string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Batch eq '%s'", batch))
	req.URL.RawQuery = params.Encode()
}
