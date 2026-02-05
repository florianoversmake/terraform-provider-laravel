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
	UpdatedAt string `json:"updated_at"`
}

type CreateRecipeRequest struct {
	Name   string `json:"name"`
	User   string `json:"user"`
	Script string `json:"script"`
}

func (c *Client) CreateRecipe(ctx context.Context, req CreateRecipeRequest) (*Recipe, error) {
	path := c.orgPath("/recipes")
	var recipe Recipe
	id, err := c.PostJsonApi(ctx, path, req, &recipe)
	if err != nil {
		return nil, err
	}
	recipe.ID = int64(id)
	return &recipe, nil
}

func (c *Client) ListRecipes(ctx context.Context) ([]Recipe, error) {
	items, err := c.GetJsonApiList(ctx, c.orgPath("/recipes"))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(r *Recipe, id int64) { r.ID = id })
}

func (c *Client) GetRecipe(ctx context.Context, recipeID int) (*Recipe, error) {
	path := c.orgPath(fmt.Sprintf("/recipes/%d", recipeID))
	var recipe Recipe
	id, err := c.GetJsonApi(ctx, path, &recipe)
	if err != nil {
		return nil, err
	}
	recipe.ID = int64(id)
	return &recipe, nil
}

func (c *Client) UpdateRecipe(ctx context.Context, recipeID int, req CreateRecipeRequest) (*Recipe, error) {
	path := c.orgPath(fmt.Sprintf("/recipes/%d", recipeID))
	var recipe Recipe
	id, err := c.PutJsonApi(ctx, path, req, &recipe)
	if err != nil {
		return nil, err
	}
	recipe.ID = int64(id)
	return &recipe, nil
}

func (c *Client) DeleteRecipe(ctx context.Context, recipeID int) error {
	path := c.orgPath(fmt.Sprintf("/recipes/%d", recipeID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type RunRecipeRequest struct {
	Servers []int64 `json:"servers"`
	Notify  bool    `json:"notify"`
}

func (c *Client) RunRecipe(ctx context.Context, recipeID int, req RunRecipeRequest) error {
	path := c.orgPath(fmt.Sprintf("/recipes/%d/runs", recipeID))
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}
