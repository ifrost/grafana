package librarypanels

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestPatchLibraryPanel(t *testing.T) {
	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel that does not exist, it should fail",
		func(t *testing.T, sc scenarioContext) {
			cmd := patchLibraryPanelCommand{}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": "unknown"})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 404, resp.Status())
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel that exists, it should succeed",
		func(t *testing.T, sc scenarioContext) {
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID, ":dashboardId": "1"})
			resp := sc.service.connectHandler(sc.reqContext)
			require.Equal(t, 200, resp.Status())
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID, ":dashboardId": "2"})
			resp = sc.service.connectHandler(sc.reqContext)
			require.Equal(t, 200, resp.Status())

			newFolder := createFolderWithACL(t, "NewFolder", sc.user, []folderACLItem{})
			cmd := patchLibraryPanelCommand{
				FolderID: newFolder.Id,
				Name:     "Panel - New name",
				Model: []byte(`
								{
								  "datasource": "${DS_GDEV-TESTDATA}",
								  "id": 1,
								  "title": "Model - New name",
								  "type": "text"
								}
							`),
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp = sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 200, resp.Status())
			var result = validateAndUnMarshalResponse(t, resp)
			var expected = libraryPanelResult{
				Result: libraryPanel{
					ID:       1,
					OrgID:    1,
					FolderID: newFolder.Id,
					UID:      sc.initialResult.Result.UID,
					Name:     "Panel - New name",
					Model: map[string]interface{}{
						"datasource": "${DS_GDEV-TESTDATA}",
						"id":         float64(1),
						"title":      "Panel - New name",
						"type":       "text",
					},
					Meta: LibraryPanelDTOMeta{
						CanEdit:             true,
						ConnectedDashboards: 2,
						Created:             sc.initialResult.Result.Meta.Created,
						Updated:             result.Result.Meta.Updated,
						CreatedBy: LibraryPanelDTOMetaUser{
							ID:        1,
							Name:      "user_in_db",
							AvatarUrl: "/avatar/402d08de060496d6b6874495fe20f5ad",
						},
						UpdatedBy: LibraryPanelDTOMetaUser{
							ID:        1,
							Name:      "signed_in_user",
							AvatarUrl: "/avatar/37524e1eb8b3e32850b57db0a19af93b",
						},
					},
				},
			}
			if diff := cmp.Diff(expected, result, getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel with folder only, it should change folder successfully and return correct result",
		func(t *testing.T, sc scenarioContext) {
			newFolder := createFolderWithACL(t, "NewFolder", sc.user, []folderACLItem{})
			cmd := patchLibraryPanelCommand{
				FolderID: newFolder.Id,
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 200, resp.Status())
			var result = validateAndUnMarshalResponse(t, resp)
			sc.initialResult.Result.FolderID = newFolder.Id
			sc.initialResult.Result.Meta.CreatedBy.Name = "user_in_db"
			sc.initialResult.Result.Meta.CreatedBy.AvatarUrl = "/avatar/402d08de060496d6b6874495fe20f5ad"
			if diff := cmp.Diff(sc.initialResult.Result, result.Result, getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel with name only, it should change name successfully, sync title and return correct result",
		func(t *testing.T, sc scenarioContext) {
			cmd := patchLibraryPanelCommand{
				FolderID: -1,
				Name:     "New Name",
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			var result = validateAndUnMarshalResponse(t, resp)
			sc.initialResult.Result.Name = "New Name"
			sc.initialResult.Result.Meta.CreatedBy.Name = "user_in_db"
			sc.initialResult.Result.Meta.CreatedBy.AvatarUrl = "/avatar/402d08de060496d6b6874495fe20f5ad"
			sc.initialResult.Result.Model["title"] = "New Name"
			if diff := cmp.Diff(sc.initialResult.Result, result.Result, getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel with model only, it should change model successfully and return correct result",
		func(t *testing.T, sc scenarioContext) {
			cmd := patchLibraryPanelCommand{
				FolderID: -1,
				Model:    []byte(`{ "title": "New Model Title", "name": "New Model Name" }`),
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			var result = validateAndUnMarshalResponse(t, resp)
			sc.initialResult.Result.Model = map[string]interface{}{
				"title": "Text - Library Panel",
				"name":  "New Model Name",
			}
			sc.initialResult.Result.Meta.CreatedBy.Name = "user_in_db"
			sc.initialResult.Result.Meta.CreatedBy.AvatarUrl = "/avatar/402d08de060496d6b6874495fe20f5ad"
			if diff := cmp.Diff(sc.initialResult.Result, result.Result, getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	scenarioWithLibraryPanel(t, "When another admin tries to patch a library panel, it should change UpdatedBy successfully and return correct result",
		func(t *testing.T, sc scenarioContext) {
			cmd := patchLibraryPanelCommand{FolderID: -1}
			sc.reqContext.UserId = 2
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			var result = validateAndUnMarshalResponse(t, resp)
			sc.initialResult.Result.Meta.UpdatedBy.ID = int64(2)
			sc.initialResult.Result.Meta.CreatedBy.Name = "user_in_db"
			sc.initialResult.Result.Meta.CreatedBy.AvatarUrl = "/avatar/402d08de060496d6b6874495fe20f5ad"
			if diff := cmp.Diff(sc.initialResult.Result, result.Result, getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel with a name that already exists, it should fail",
		func(t *testing.T, sc scenarioContext) {
			command := getCreateCommand(sc.folder.Id, "Another Panel")
			resp := sc.service.createHandler(sc.reqContext, command)
			var result = validateAndUnMarshalResponse(t, resp)
			cmd := patchLibraryPanelCommand{
				Name: "Text - Library Panel",
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": result.Result.UID})
			resp = sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 400, resp.Status())
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel with a folder where a library panel with the same name already exists, it should fail",
		func(t *testing.T, sc scenarioContext) {
			newFolder := createFolderWithACL(t, "NewFolder", sc.user, []folderACLItem{})
			command := getCreateCommand(newFolder.Id, "Text - Library Panel")
			resp := sc.service.createHandler(sc.reqContext, command)
			var result = validateAndUnMarshalResponse(t, resp)
			cmd := patchLibraryPanelCommand{
				FolderID: 1,
			}
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": result.Result.UID})
			resp = sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 400, resp.Status())
		})

	scenarioWithLibraryPanel(t, "When an admin tries to patch a library panel in another org, it should fail",
		func(t *testing.T, sc scenarioContext) {
			cmd := patchLibraryPanelCommand{
				FolderID: sc.folder.Id,
			}
			sc.reqContext.OrgId = 2
			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
			resp := sc.service.patchHandler(sc.reqContext, cmd)
			require.Equal(t, 404, resp.Status())
		})
}
