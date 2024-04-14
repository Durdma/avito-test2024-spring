package httpv1

import (
	"avito-test2024-spring/internal/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

//func initTestServer(cfg *config.Config) {
//	logs := logger.NewLogs(cfg.Logger)
//
//	logs.Logger.Info().Msg("Starting Tests app")
//	logs.Logger.Info().Interface("config", cfg).Msg("")
//
//	cache := cache2.NewRedisCache(cfg.Cache)
//	logs.Logger.Info().Msg("Initialized connection pool Cache")
//
//	dbPool := postgresql.NewConnectionPool(cfg.PostgreSQL, logs)
//	logs.Logger.Info().Msg("Initialized connection pool DB")
//
//	repos := repository.NewRepositories(dbPool)
//	logs.Logger.Info().Msg("Initialized repos")
//
//	tokenManager, err := auth.NewManager(cfg.JWT.SigningKey)
//	if err != nil {
//		logs.Logger.Error().Err(err).Msg("error occured while init of token manager")
//	}
//	logs.Logger.Info().Msg("Initialized tokenManager")
//
//	services := service.NewServices(repos, tokenManager, cache)
//	logs.Logger.Info().Msg("Initialized services")
//
//	handlers := controller.NewHandler(services.Banners, services.Tags, services.Features, services.Users, logs, tokenManager, cache)
//	logs.Logger.Info().Msg("Initialized handlers")
//
//	srv := server.NewServer(cfg.HTTP, handlers.Init("localhost", cfg.HTTP.Port))
//	go func() {
//		if err := srv.Run(); err != nil {
//			logs.Logger.Error().Err(err).Msg("error occurred while running http server")
//		}
//	}()
//
//	logs.Logger.Info().Msg("Test server started")
//}

func TestUserBanner(t *testing.T) {
	//cfg, err := config.Init("..\\..\\configs\\main.yaml")
	//if err != nil {
	//	log.Fatal(err.Error())
	//}

	//initTestServer(cfg)

	client := &http.Client{}

	r1Body := make(map[string]interface{})
	r1Body["is_admin"] = true
	r1BodyJSON, err := json.Marshal(r1Body)
	require.NoError(t, err)

	r1, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v:%v/api/v1/users", "localhost", "8080"), strings.NewReader(string(r1BodyJSON)))
	r1.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r1)
	require.NoError(t, err)

	if resp.StatusCode != 201 || resp.Body == nil {
		log.Fatal(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	accessToken := make(map[string]string)
	if err := json.Unmarshal(body, &accessToken); err != nil {
		fmt.Println("r1")
		log.Fatal(err)
	}

	r2, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v:%v/api/v1/tags", "localhost", "8080"), nil)
	r2.Header.Add("Content-Type", "application/json")
	r2.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r2)
	require.NoError(t, err)

	if resp.StatusCode != 201 || resp.Body == nil {
		log.Fatal(resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	tagId := make(map[string]int)
	if err := json.Unmarshal(body, &tagId); err != nil {
		log.Fatal(err, "r2")
	}

	r3, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v:%v/api/v1/features", "localhost", "8080"), nil)
	r3.Header.Add("Content-Type", "application/json")
	r3.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r3)
	require.NoError(t, err)

	if resp.StatusCode != 201 || resp.Body == nil {
		log.Fatal(resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	FeatureId := make(map[string]int)
	if err := json.Unmarshal(body, &FeatureId); err != nil {
		log.Fatal(err, "r3")
	}

	r4Body := make(map[string]interface{})
	r4Body["is_admin"] = false
	r4Body["tag_id"] = tagId["tag_id"]
	r4BodyJSON, err := json.Marshal(r4Body)
	require.NoError(t, err)

	r4, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v:%v/api/v1/users", "localhost", "8080"), strings.NewReader(string(r4BodyJSON)))
	r4.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(r4)
	require.NoError(t, err)

	if resp.StatusCode != 201 || resp.Body == nil {
		fmt.Println("r4")
		log.Fatal(resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	accessTokenUser := make(map[string]string)
	if err := json.Unmarshal(body, &accessTokenUser); err != nil {
		log.Fatal(err, "r4")
	}

	r5Body := make(map[string]interface{})
	r5Body["tags_ids"] = []int{tagId["tag_id"]}
	r5Body["feature_id"] = FeatureId["feature_id"]
	r5Body["content"] = map[string]string{"title": "some title", "text": "some_text", "url": "http://example.com"}
	r5Body["is_active"] = true
	r5BodyJSON, err := json.Marshal(r5Body)
	require.NoError(t, err)

	r5, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v:%v/api/v1/banner", "localhost", "8080"), strings.NewReader(string(r5BodyJSON)))
	r5.Header.Add("Content-Type", "application/json")
	r5.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r5)
	require.NoError(t, err)

	if resp.StatusCode != 201 || resp.Body == nil {
		fmt.Println("r5")
		log.Fatal(resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	bannerId := make(map[string]int)
	if err := json.Unmarshal(body, &bannerId); err != nil {
		log.Fatal(err, "r5")
	}

	fmt.Println(bannerId)
	fmt.Println("no errors")

	t.Run("Successful_GetUserBanner", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("http://%v:%v/api/v1/user_banner?tag_id=%v&feature_id=%v", "localhost", "8080", tagId["tag_id"], FeatureId["feature_id"]), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessTokenUser["access_token"]))

		resp, err = client.Do(request)
		require.NoError(t, err)

		if resp.StatusCode != 200 {
			log.Fatal(resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		banner := models.Banner{}

		err = json.Unmarshal(body, &banner)
		require.NoError(t, err)

		log.Printf("banner with tag_id=%v and feature_id=%v\n%v\n", tagId["tag_id"], FeatureId["feature_id"], banner)
	})

	r6, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://%v:%v/api/v1/banner/%v", "localhost", "8080", bannerId["banner_id"]), nil)
	r6.Header.Add("Content-Type", "application/json")
	r6.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r6)
	require.NoError(t, err)

	log.Println(tagId)

	r7, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://%v:%v/api/v1/tags/%v", "localhost", "8080", tagId["tag_id"]), nil)
	r7.Header.Add("Content-Type", "application/json")
	r7.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r7)
	require.NoError(t, err)

	r8, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://%v:%v/api/v1/features/%v", "localhost", "8080", FeatureId["feature_id"]), nil)
	r8.Header.Add("Content-Type", "application/json")
	r8.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken["access_token"]))

	resp, err = client.Do(r8)
	require.NoError(t, err)
}
