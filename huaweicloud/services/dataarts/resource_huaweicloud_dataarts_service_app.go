// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArtsStudio
// ---------------------------------------------------------------

package dataarts

import (
	"context"
	"encoding/json"
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

func ResourceServiceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAppCreate,
		UpdateContext: resourceServiceAppUpdate,
		ReadContext:   resourceServiceAppRead,
		DeleteContext: resourceServiceAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAppImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of DLM.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the app.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the app.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The type of the app.`,
			},
			"app_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the app.`,
			},
			"app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secret of the app.`,
			},
		},
	}
}

func resourceServiceAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createApp: create an app
	var (
		createAppHttpUrl = "v1/{project_id}/service/apps"
		createAppProduct = "dataarts"
	)
	createAppClient, err := cfg.NewServiceClient(createAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createAppPath := createAppClient.Endpoint + createAppHttpUrl
	createAppPath = strings.ReplaceAll(createAppPath, "{project_id}", createAppClient.ProjectID)

	createAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm_type":     d.Get("dlm_type").(string),
		},
	}

	createAppOpt.JSONBody = utils.RemoveNil(buildCreateAppBodyParams(d))
	createAppResp, err := createAppClient.Request("POST", createAppPath, &createAppOpt)
	if err != nil {
		return diag.Errorf("error creating app: %s", err)
	}

	createAppRespBody, err := utils.FlattenResponse(createAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("id", createAppRespBody)
	if err != nil {
		return diag.Errorf("error creating app: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceServiceAppRead(ctx, d, meta)
}

func buildCreateAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
		"app_type":    utils.ValueIngoreEmpty(d.Get("app_type")),
	}
	return bodyParams
}

func resourceServiceAppRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getApp: Query the app
	var (
		getAppHttpUrl = "v1/{project_id}/service/apps/{id}"
		getAppProduct = "dataarts"
	)
	getAppClient, err := cfg.NewServiceClient(getAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getAppPath := getAppClient.Endpoint + getAppHttpUrl
	getAppPath = strings.ReplaceAll(getAppPath, "{project_id}", getAppClient.ProjectID)
	getAppPath = strings.ReplaceAll(getAppPath, "{id}", d.Id())

	getAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm_type":     d.Get("dlm_type").(string),
		},
	}

	getAppResp, err := getAppClient.Request("GET", getAppPath, &getAppOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, parseAppNotFoundError(err), "error retrieving app")
	}

	getAppRespBody, err := utils.FlattenResponse(getAppResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getAppRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getAppRespBody, nil)),
		d.Set("app_type", utils.PathSearch("app_type", getAppRespBody, nil)),
		d.Set("app_key", utils.PathSearch("app_key", getAppRespBody, nil)),
		d.Set("app_secret", utils.PathSearch("app_secret", getAppRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceServiceAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAppChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateAppChanges...) {
		// updateApp: update the App
		var (
			updateAppHttpUrl = "v1/{project_id}/service/apps/{id}"
			updateAppProduct = "dataarts"
		)
		updateAppClient, err := cfg.NewServiceClient(updateAppProduct, region)
		if err != nil {
			return diag.Errorf("error creating DataArts Studio client: %s", err)
		}

		updateAppPath := updateAppClient.Endpoint + updateAppHttpUrl
		updateAppPath = strings.ReplaceAll(updateAppPath, "{project_id}", updateAppClient.ProjectID)
		updateAppPath = strings.ReplaceAll(updateAppPath, "{id}", d.Id())

		updateAppOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"workspace":    d.Get("workspace_id").(string),
				"dlm_type":     d.Get("dlm_type").(string),
			},
		}

		updateAppOpt.JSONBody = utils.RemoveNil(buildUpdateAppBodyParams(d))
		_, err = updateAppClient.Request("PUT", updateAppPath, &updateAppOpt)
		if err != nil {
			return diag.Errorf("error updating app: %s", err)
		}
	}
	return resourceServiceAppRead(ctx, d, meta)
}

func buildUpdateAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIngoreEmpty(d.Get("name")),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceServiceAppDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteApp: delete the app
	var (
		deleteAppHttpUrl = "v1/{project_id}/service/apps/{id}"
		deleteAppProduct = "dataarts"
	)
	deleteAppClient, err := cfg.NewServiceClient(deleteAppProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deleteAppPath := deleteAppClient.Endpoint + deleteAppHttpUrl
	deleteAppPath = strings.ReplaceAll(deleteAppPath, "{project_id}", deleteAppClient.ProjectID)
	deleteAppPath = strings.ReplaceAll(deleteAppPath, "{id}", d.Id())

	deleteAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm_type":     d.Get("dlm_type").(string),
		},
	}

	_, err = deleteAppClient.Request("DELETE", deleteAppPath, &deleteAppOpt)
	if err != nil {
		return diag.Errorf("error deleting app: %s", err)
	}

	return nil
}

func resourceAppImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<dlm_type>/<id>")
	}

	d.Set("workspace_id", parts[0])
	d.Set("dlm_type", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func parseAppNotFoundError(respErr error) error {
	var apiErr interface{}
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr != nil {
			return pErr
		}
		errCode, err := jmespath.Search(`error_code`, apiErr)
		if err != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", err.Error())
		}

		if errCode == `DLM.4063` {
			return golangsdk.ErrDefault404{}
		}
	}
	return respErr
}
