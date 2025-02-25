package accesstoartifacts

import (
	"testing"

	"github.com/aquasecurity/chain-bench/internal/checks/common"
	"github.com/aquasecurity/chain-bench/internal/checks/consts"
	"github.com/aquasecurity/chain-bench/internal/models/checkmodels"
	"github.com/aquasecurity/chain-bench/internal/testutils"
	"github.com/aquasecurity/chain-bench/internal/testutils/builders"
)

const (
	vulnerabilityScanningTask = "argonsecurity/scanner-action"
)

func TestAccessToArtifactsChecker(t *testing.T) {
	tests := []testutils.CheckTest{
		{
			Name: "no org settings permissions",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithOrganization(builders.NewOrganizationBuilder().Build()).
					WithPackageRegistry(builders.NewRegistryBuilder().WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Unknown, Details: consts.Details_organization_missingMinimalPermissions}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "no org packages permissions",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Unknown, Details: consts.Details_organization_hooks_missingMinimalPermissions}),
			},
		},
		{
			Name: "no org & no repo permissions",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "repo permissions only, no org permissions",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "Package registry with 2mfa disabled",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithOrganization(builders.NewOrganizationBuilder().WithReposDefaultPermissions("read").Build()).
					WithPackageRegistry(builders.NewRegistryBuilder().WithTwoFactorAuthenticationEnabled(false).WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "Package registry with 2mfa enabled",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithOrganization(builders.NewOrganizationBuilder().WithReposDefaultPermissions("read").Build()).
					WithPackageRegistry(builders.NewRegistryBuilder().WithTwoFactorAuthenticationEnabled(true).WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Passed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "Package registry with 2 public packages under private repo",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().WithTwoFactorAuthenticationEnabled(false).WithPackages("npm", "public", true).WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "2 anonymous accessed packages"}),
			},
		},
		{
			Name: "Package registry with 1 private and 1 public packages under private repo",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().WithTwoFactorAuthenticationEnabled(false).WithPackages("npm", "private", true).WithPackages("npm", "public", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed, Details: "1 anonymous accessed packages"}),
			},
		},
		{
			Name: "Package registry with 2 private packages",
			Data: &checkmodels.CheckData{
				AssetsMetadata: builders.NewAssetsDataBuilder().
					WithPackageRegistry(builders.NewRegistryBuilder().WithTwoFactorAuthenticationEnabled(false).WithPackages("npm", "private", true).WithPackages("npm", "private", true).Build()).
					Build(),
			},
			Expected: []*checkmodels.CheckRunResult{
				checkmodels.ToCheckRunResult("4.2.3", checksMetadata.Checks["4.2.3"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Failed}),
				checkmodels.ToCheckRunResult("4.2.5", checksMetadata.Checks["4.2.5"], checksMetadata.Url, &checkmodels.CheckResult{Status: checkmodels.Passed}),
			},
		},
	}
	testutils.RunCheckTests(t, common.GetRegoRunAction(regoQuery, checksMetadata), tests)
}
