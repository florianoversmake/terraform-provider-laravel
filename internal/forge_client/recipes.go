package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Recipe struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Script    string `json:"script"`
	CreatedAt string `json:"created_at"`
}

type recipeResponse struct {
	Recipe Recipe `json:"recipe"`
}

type recipesResponse struct {
	Recipes []Recipe `json:"recipes"`
}

type CreateRecipeRequest struct {
	Name   string `json:"name"`
	User   string `json:"user"`
	Script string `json:"script"`
}

func (c *Client) CreateRecipe(ctx context.Context, req CreateRecipeRequest) (*Recipe, error) {
	path := "/recipes"
	var res recipeResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Recipe, nil
}

func (c *Client) ListRecipes(ctx context.Context) ([]Recipe, error) {
	path := "/recipes"
	var res recipesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Recipes, nil
}

func (c *Client) GetRecipe(ctx context.Context, recipeID int) (*Recipe, error) {
	path := fmt.Sprintf("/recipes/%d", recipeID)
	var res recipeResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Recipe, nil
}

func (c *Client) UpdateRecipe(ctx context.Context, recipeID int, req CreateRecipeRequest) (*Recipe, error) {
	path := fmt.Sprintf("/recipes/%d", recipeID)
	var res recipeResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Recipe, nil
}

func (c *Client) DeleteRecipe(ctx context.Context, recipeID int) error {
	path := fmt.Sprintf("/recipes/%d", recipeID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type RunRecipeRequest struct {
	Servers []int64 `json:"servers"`
	Notify  bool    `json:"notify"`
}

func (c *Client) RunRecipe(ctx context.Context, recipeID int, req RunRecipeRequest) error {
	path := fmt.Sprintf("/recipes/%d/run", recipeID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}
