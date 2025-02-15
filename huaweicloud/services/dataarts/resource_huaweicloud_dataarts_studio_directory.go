package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// API: DataArtsStudio POST /v2/{project_id}/design/directorys
// API: DataArtsStudio DELETE /v2/{project_id}/design/directorys
// API: DataArtsStudio GET /v2/{project_id}/design/directorys
// API: DataArtsStudio PUT /v2/{project_id}/design/directorys
func ResourceDataArtsStudioDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudioDirectoryCreate,
		ReadContext:   resourceStudioDirectoryRead,
		UpdateContext: resourceStudioDirectoryUpdate,
		DeleteContext: resourceStudioDirectoryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceStudioDirectoryImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"qualified_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStudioDirectoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	//nolint:misspell
	createDirectoryHttpUrl := "v2/{project_id}/design/directorys"
	createDirectoryProduct := "dataarts"

	createDirectoryClient, err := cfg.NewServiceClient(createDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V2 Client: %s", err)
	}
	createDirectoryPath := createDirectoryClient.Endpoint + createDirectoryHttpUrl
	createDirectoryPath = strings.ReplaceAll(createDirectoryPath, "{project_id}", createDirectoryClient.ProjectID)

	createDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createDirectoryOpt.JSONBody = utils.RemoveNil(buildCreateDirectoryBodyParams(d))
	createDirectoryResp, err := createDirectoryClient.Request("POST", createDirectoryPath, &createDirectoryOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	createDirectoryRespBody, err := utils.FlattenResponse(createDirectoryResp)
	if err != nil {
		return diag.FromErr(err)
	}
	directory := utils.PathSearch("data.value", createDirectoryRespBody, nil)

	id, err := jmespath.Search("id", directory)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio directory: %s is not found in API response", "id")
	}

	// need to set qualified name to filter result in READ.
	qualifiedName, err := jmespath.Search("qualified_name", directory)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio directory: %s is not found in API response", "qualifiedName")
	}
	d.SetId(id.(string))
	d.Set("qualified_name", qualifiedName)

	return resourceStudioDirectoryRead(ctx, d, meta)
}

func buildCreateDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
		"parent_id":   utils.ValueIngoreEmpty(d.Get("parent_id")),
	}
	return bodyParams
}

func resourceStudioDirectoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	//nolint:misspell
	getDirectoryHttpUrl := "v2/{project_id}/design/directorys?type={type}"
	getDirectoryProduct := "dataarts"

	getDirectoryClient, err := cfg.NewServiceClient(getDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V2 Client: %s", err)
	}

	getDirectoryPath := getDirectoryClient.Endpoint + getDirectoryHttpUrl
	getDirectoryPath = strings.ReplaceAll(getDirectoryPath, "{project_id}", getDirectoryClient.ProjectID)
	getDirectoryPath = strings.ReplaceAll(getDirectoryPath, "{type}", d.Get("type").(string))

	getDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getDirectoryResp, err := getDirectoryClient.Request("GET", getDirectoryPath, &getDirectoryOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	getDirectoryRespBody, err := utils.FlattenResponse(getDirectoryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	paths := strings.Split(d.Get("qualified_name").(string), ".")
	jsonPaths := fmt.Sprintf("data.value[?name=='%s']", paths[0])
	for i, path := range paths {
		if i == 0 {
			continue
		}
		jsonPaths += fmt.Sprintf("[children][][?name=='%s'][]", path)
	}

	directories := utils.PathSearch(jsonPaths, getDirectoryRespBody, make([]interface{}, 0)).([]interface{})
	if len(directories) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DataArts Studio directory")
	}

	directory := directories[0]
	d.SetId(utils.PathSearch("id", directory, "").(string))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("name", utils.PathSearch("name", directory, nil)),
		d.Set("type", utils.PathSearch("type", directory, nil)),
		d.Set("description", utils.PathSearch("description", directory, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", directory, nil)),
		d.Set("root_id", utils.PathSearch("root_id", directory, nil)),
		d.Set("qualified_name", utils.PathSearch("qualified_name", directory, nil)),
		d.Set("created_at", utils.PathSearch("create_time", directory, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", directory, nil)),
		d.Set("created_by", utils.PathSearch("create_by", directory, nil)),
		d.Set("updated_by", utils.PathSearch("update_by", directory, nil)),
		d.Set("children", utils.PathSearch(`children[*].name`, directory, make([]interface{}, 0))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DataArts Studio directory fields: %s", err)
	}

	return nil
}

func resourceStudioDirectoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	//nolint:misspell
	updateDirectoryHttpUrl := "v2/{project_id}/design/directorys"
	updateDirectoryProduct := "dataarts"

	updateDirectoryClient, err := cfg.NewServiceClient(updateDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V2 Client: %s", err)
	}
	updateDirectoryPath := updateDirectoryClient.Endpoint + updateDirectoryHttpUrl
	updateDirectoryPath = strings.ReplaceAll(updateDirectoryPath, "{project_id}", updateDirectoryClient.ProjectID)

	updateDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	updateDirectoryOpt.JSONBody = utils.RemoveNil(buildUpdateDirectoryBodyParams(d))
	updateDirectoryResp, err := updateDirectoryClient.Request("PUT", updateDirectoryPath, &updateDirectoryOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	updateDirectoryRespBody, err := utils.FlattenResponse(updateDirectoryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	directory := utils.PathSearch("data.value", updateDirectoryRespBody, nil)

	// if you change the parent id, the qualified name will be changed, need to set to filter result in READ.
	qualifiedName, err := jmespath.Search("qualified_name", directory)
	if err != nil {
		return diag.Errorf("error updating DataArts Studio directory: %s is not found in API response", "qualifiedName")
	}
	if qualifiedName == nil {
		qualifiedName = d.Get("qualified_name")
	}
	d.Set("qualified_name", qualifiedName)

	return resourceStudioDirectoryRead(ctx, d, meta)
}

func buildUpdateDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":          d.Id(),
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
		"parent_id":   utils.ValueIngoreEmpty(d.Get("parent_id")),
	}
	return bodyParams
}

func resourceStudioDirectoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	//nolint:misspell
	deleteDirectoryHttpUrl := "v2/{project_id}/design/directorys?ids={id}"
	deleteDirectoryProduct := "dataarts"

	deleteDirectoryClient, err := cfg.NewServiceClient(deleteDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio V2 Client: %s", err)
	}
	deleteDirectoryPath := deleteDirectoryClient.Endpoint + deleteDirectoryHttpUrl
	deleteDirectoryPath = strings.ReplaceAll(deleteDirectoryPath, "{project_id}", deleteDirectoryClient.ProjectID)
	deleteDirectoryPath = strings.ReplaceAll(deleteDirectoryPath, "{id}", d.Id())

	deleteDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	_, err = deleteDirectoryClient.Request("DELETE", deleteDirectoryPath, &deleteDirectoryOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceStudioDirectoryImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<type>/<qualified_name>")
	}

	d.Set("workspace_id", parts[0])
	d.Set("type", parts[1])
	d.Set("qualified_name", parts[2])

	return []*schema.ResourceData{d}, nil
}
