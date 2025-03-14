package api

import (
	"WaterMark/src/cmd"
	"WaterMark/src/images"
	"strings"
)

// getTplListType
//
//	@return map
func getTplListType() map[string]string {
	r := map[string]string{}
	for _, v := range images.GetTemplates() {
		r[v.ID] = v.Type
	}
	return r
}

// getImageWaterMarkPreview
//
//	@param tid
//	@param path
//	@param color
//	@param onlyBottomBorder
//	@return map
func getImageWaterMarkPreview(tid string, path string, color string, onlyBottomBorder bool) map[string]string {
	e := images.NewExternal()
	e.WithBoderColor(color).WithOnlyBottomFlag(onlyBottomBorder).WithPath(path).WithTid(tid)

	r := images.GetPreviewWaterMark(e)
	r["SaveImgPath"] = cmd.GetPwdPath(strings.TrimLeft(r["SaveImgPath"], "."))
	return r
}
