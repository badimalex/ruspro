package server

import (
	"context"
	"fmt"
	"net/http"
	pb "ruspro/api"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Server struct {
	pb.UnimplementedRusProfileServiceServer
}

func (s *Server) GetCompanyInfo(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyResponse, error) {
	if req.Inn == "" {
		return nil, fmt.Errorf("INN is required")
	}

	response, err := queryRusProfile(req.Inn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func queryRusProfile(inn string) (*pb.CompanyResponse, error) {
	url := "https://www.rusprofile.ru/search?query=" + inn + "&type=ul"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	kpp := doc.Find("#clip_kpp").Text()

	companyName := doc.Find(".company-name").Text()

	ceoName := doc.Find(".company-info__text a.link-arrow").Text()

	if companyName == "" || kpp == "" || ceoName == "" {
		return nil, fmt.Errorf("information not found")
	}

	return &pb.CompanyResponse{
		Inn:  strings.TrimSpace(inn),
		Kpp:  strings.TrimSpace(kpp),
		Name: strings.TrimSpace(companyName),
		Ceo:  strings.TrimSpace(ceoName),
	}, nil
}
