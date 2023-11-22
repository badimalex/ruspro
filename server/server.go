package server

import (
	"context"
	"fmt"
	"net/http"
	pb "ruspro/api"
	"ruspro/internal/logging"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Server struct {
	ApiUrl string
	pb.UnimplementedRusProfileServiceServer
}

func (s *Server) GetCompanyInfo(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyResponse, error) {
	logging.Log.Info("Start fetching GetCompanyInfo")

	if req.Inn == "" {
		logging.Log.Error("failed to fetching INN")
		return nil, fmt.Errorf("INN is required")
	}

	response, err := queryRusProfile(req.Inn, s.ApiUrl)
	if err != nil {
		logging.Log.Error("failed to query Profile")
		return nil, err
	}

	return response, nil
}

func queryRusProfile(inn, apiUrl string) (*pb.CompanyResponse, error) {
	if !isValidINN(inn) {
		logging.Log.Error("INN is not valid")
		return nil, fmt.Errorf("некорректный ИНН")
	}

	url := strings.ReplaceAll(apiUrl, "{{inn}}", inn)
	logging.Log.Infof("query fetching rusprofile url %s", url)

	resp, err := http.Get(url)
	if err != nil {
		logging.Log.Error("error query to server rusprofile")
		return nil, fmt.Errorf("ошибка запроса к RusProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Log.Error("api rusprofile is not working")
		return nil, fmt.Errorf("сервис RusProfile недоступен, код ответа: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logging.Log.Error("error parse body")
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	kpp := doc.Find("#clip_kpp").Text()
	companyName := doc.Find(".company-name").Text()
	ceoName := doc.Find(".company-info__text a.link-arrow").Text()

	if companyName == "" || kpp == "" || ceoName == "" {
		logging.Log.Error("company not found")
		return nil, fmt.Errorf("компания не найдена")
	}

	logging.Log.Info("Success company fetching")

	return &pb.CompanyResponse{
		Inn:  strings.TrimSpace(inn),
		Kpp:  strings.TrimSpace(kpp),
		Name: strings.TrimSpace(companyName),
		Ceo:  strings.TrimSpace(ceoName),
	}, nil
}
