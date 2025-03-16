// 显示首页
function showHome() {
    window.location.href = "./index.html"
}

// 显示快门查询
function showShutterimg() {
    window.location.href = "./shutter.html"
}

// 显示图片导出
function showWaterMark() {
    window.location.href = "./watermark.html"
}

// 显示水印生成
function showWaterMarkProcess() {
    window.location.href = "./watermarkExport.html"
}

function sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}

// 获取请求url
function getReqUrl(type, params) {
    var host = SERVER_HOST
    var url = "server/getImagePreview"
    switch (type) {
        case "shutter":
            url = "server/getShutterByFile"
            break
        case "ImagePreview":
            url = "server/getImagePreview"
            break
        case "ImageResize":
            url = "server/addImageResizeTask"
            break
        case "ImagePreviewSmall":
            url = "server/imagePreviewSmall"
            break
        case "TplListType":
            url = "server/getTplListType"
            break
        case "WaterMarkPreview":
            url = "server/getImageWaterMarkPreview"
            break
        case "AddPreviewTask":
            url = "server/addPreviewTask"
            break
        case "Download":
            url = "server/downloadFile"
            break
        case "ChangeImagePath":
            url = "server/changeImagePath"
            break
    }
    url = url + "?random="+Math.random()
    if (params.length > 0) {
        url = url + "&" + params.join("&")
    }
    return host + url
}

// 选择图片文件
function shutterSelectImage() {
    window.go.gui.App.SelectImageFile().then(result => {
        $("#shutterimg-input").val(result)
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}

// 查询快门信息接口
function shutterImgUpload() {
    $("#shutterimg-result").html("")
    $("#shutterimg-preview-img").hide()
    $("#shutterimg-preview-img").attr("src", "")
    
    var file = $("#shutterimg-input").val()
    if (file.length == 0) {
        $("#shutterimg-toast").text("请选择需要查看的文件!").fadeIn(400).delay(500).fadeOut(400); 
        return false
    }
    getExifInfo(file)
    $("#shutterimg-preview-img").attr("src", getReqUrl("ImagePreview", ["imagePath="+file]))
    $("#shutterimg-preview-img").show()
    return true
}

// 获取exif 信息
function getExifInfo(file) {
    var result
    var data = new FormData();
    data.append('shutterimg', file);
    $.ajax({
        url : getReqUrl("shutter", []),
        type : "post",
        data : data,
        cache : false,
        async : false,
        processData : false,
        contentType : false,
        success : function (response) {
            if (response['ShutterCount'] == 0 && response['MechanicalShutterCount'] == 0) {
                $("#shutterimg-toast").text("查看失败! 请确认选择的图片保存有exif信息").fadeIn(400).delay(500).fadeOut(400); 
            } else {
                var li = "机器快门次数:"+response['MechanicalShutterCount']+" 快门次数:"+response['ShutterCount']
                $("#shutterimg-result").append(li)
            }
            result = response
        },
        error: function(xhr, status, error) {
            $("#shutterimg-toast").text("查看失败! 原因是:" + error).fadeIn(400).delay(500).fadeOut(400); 
        }
    });
    return result
}

// 选择文件夹
var g_type;
function SelectDirectory(type) {
    g_type = type
    window.go.gui.App.SelectDirectory().then(result => {
        var id = "watermarkprocess-span-" + g_type
        $("#" + id).val(result)
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}

// 批量选择图片文件
function watermarkOpenMultipleFilesDialog() {
    window.go.gui.App.SelectMultipleImageFile().then(async result => {
        if (result.length > 0) {
            var list = result.split(",");
            var limit = []
            if (list.length > 50) {
                for (var i = 0; i < 50; i++) {
                    limit[i] = list[i]
                }
            } else {
                limit = list
            }
            result = limit.join(",")
            $(".watermark-tpl-list").hide()
            $(".main-wrap").show()
            await sleep(100) // 休眠一定时间展示loading动画

            asynchronousPreviewTask(result)
            $("#div-selectImages").hide()
            $("#watermarkOpenMultipleFiles").val(result)
            
            // 异步添加图片裁剪,防止多图预览的时候页面崩溃
            addImageResizeTask(limit)
            // 加载预览图
            loadPreviewImage(limit[0], {})
            // 加载模板选项
            loadSelectTemplate()
            for (var i = 0; i < limit.length; i ++) {
                var url = getReqUrl("ImagePreviewSmall", ["imagePath="+limit[i]])
                if (i == 0) {
                    var img = "<img class='img-list img-list-selected pointer' data-src='"+ limit[i] +"' src='" + url + "' onclick='changePreviewImage(this)'>"
                } else {
                    var img = "<img class='img-list pointer lozad' data-src='"+ limit[i] +"' src='"+ url +"' onclick='changePreviewImage(this)'>"
                }
                $("#watermarkShowMultipleFiles").append(img)
                await sleep(100);
            }
            removeLoading()
        }
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}

// 移除loading
function removeLoading() {
    $(".main-wrap").hide()
    $(".watermark-tpl-list").show()
}

// 添加异步图片裁剪
function addImageResizeTask(list) {
    // 根据选择的文件异步的创建图片预览任务加速图片展示
    var data = new FormData();
    data.append("images", list.join(","))
    $.ajax({
        url : getReqUrl("ImageResize", []),
        type : "post",
        data : data,
        cache : false,
        async: false,
        processData : false,
        contentType : false,
        success : function (response) {},
        error: function(xhr, status, error) {}
    });
}

// 添加异步任务
function asynchronousPreviewTask(imageStr) {
    // 根据选择的文件异步的创建图片预览任务加速图片展示
    // 任务跳过第一张图片
    var list = imageStr.split(",")
    if (list.length <= 1) {
        return
    }
    var images = []
    for (var i = 1; i < list.length; i++) {
        images[i - 1] = list[i]
    }
    var data = new FormData();
    data.append("images", images.join(","))
    $.ajax({
        url : getReqUrl("AddPreviewTask", []),
        type : "post",
        data : data,
        cache : false,
        async: false,
        processData : false,
        contentType : false,
        success : function (response) {},
        error: function(xhr, status, error) {}
    });
}

// 加载图片模板类型
function loadSelectTemplate() {
    var data = new FormData();
    $.ajax({
        url : getReqUrl("TplListType", []),
        type : "get",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            for (var i in response) {
                var li = "<option value="+ i +">"+ response[i] +"</option>"
                $("#input-Template").append(li)   
            }
        },
        error: function(xhr, status, error) {
            
        }
    });
}

// 加载实际生成的预览图
function loadPreviewImage(file, params) {
    var data = new FormData();
    isEmpty = Object.keys(params).length === 0;
    if (isEmpty) {
        var defaultTid = "1"
        var defaultTBorderColor = "255,255,255,255"
        data.append("tid", defaultTid)
        data.append("borderColor", defaultTBorderColor)
        data.append("flag", "false")
        data.append("type", "create")
    } else {
        for (var i in params) {
            if (i == "flag") {
                data.append("flag", params[i].toString())
            } else {
                data.append(i, params[i])
            }
        }
    }
    data.append("imagePath", file)
    $.ajax({
        url : getReqUrl("WaterMarkPreview", []),
        type : "POST",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            var url = getReqUrl("ImagePreview", ["imagePath="+response["SaveImgPath"]])
            // 预览图片
            var imgContainer = "<img class='img-imagesContainer' id='img-imagesContainer' src='"+ url+"'>"
            // 添加预览图片
            $("#div-imagesContainer").html("").append(imgContainer)
            // 预览容器显示
            $("#div-templateContainer").show()
            // 设置预览参数
            $("#input-Color").val(response["BorderColors"])
            $("#input-FirstBorderColor").val(response["FirstBorderColor"])
            $("#input-SecondBorderColor").val(response["SecondBorderColor"])
            // 保存预览的源文件
            $("#input-PreviewSourceImageFile").val(file)
            // 保存预览的目标文件
            $("#input-PreviewImageFile").val(response["SaveImgPath"])

            if (response["model"] === undefined) {
                setTemplateInfo(getExifInfo(file))
            } else {
                // 如果是load,需要加载
                $("#input-Model").val(response["model"])
                $("#input-LensModel").val(response["lensModel"])
                $("#input-Params").val(response["words"])
            }
        },
        error: function(xhr, status, error) {
            
        }
    });
}

// 填充图片预览信息
function setTemplateInfo(exifInfo) {
    var param = exifInfo.FocalLength + " " + exifInfo.FNumberStr + " " + exifInfo.ExposureTime + " " + exifInfo.ISOStr
    $("#input-Model").val(exifInfo.Model)
    $("#input-LensModel").val(exifInfo.LensModel)
    $("#input-Params").val(param)
}

// 手动点击预览图片
function waterMarkPreivew() {
    var File = $("#input-PreviewSourceImageFile").val()
    var Tid = $("#input-Template").val()
    var Model = $("#input-Model").val()
    var LensModel = $("#input-LensModel").val()
    var Param = $("#input-Params").val()
    var Color = $("#input-Color").val()
    var FirstBorderColor = $("#input-FirstBorderColor").val()
    var SecondBorderColor = $("#input-SecondBorderColor").val()
    var OnlyBottomFlag = $("#input-OnlyBottomBorder:checked").val() === "on"

    loadPreviewImage(File, {
        "tid": Tid, 
        "borderColor": Color, 
        "flag" : OnlyBottomFlag,
        "words": Param,
        "model": Model,
        "lensModel": LensModel,
        "firstWordsColor": FirstBorderColor,
        "secondBorderColor": SecondBorderColor,
        "type": "create"
    })
}

// 图片导出
function waterMarkExport() {
    var File = $("#input-PreviewSourceImageFile").val()
    var PreviewFile = $("#input-PreviewImageFile").val()
    var data = new FormData
    data.append("source", File)
    data.append("preview", PreviewFile)
    $.ajax({
        url : getReqUrl("Download", []),
        type : "POST",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            $(".toast-body").html(response["path"])
            $('.toast').toast("show")
        },
        error: function(xhr, status, error) {
            
        }
    });
}

// 切换预览图片
function changePreviewImage(obj) {
    // 重置模板选择
    $("#input-Template").prop("selectedIndex", 0);
    // 获取新的预览文件
    var filePath = $(obj).attr("data-src")
    // 更新预览区域
    var tid = $("#input-Template").val()
    console.log(tid)
    loadPreviewImage(filePath, {
        "tid": tid, 
        "flag" : false,
        "type": "load"
    })
    $(".img-list").each(function() {
        $(this).removeClass("img-list-selected")
    });
    $(obj).addClass("img-list-selected")
}

// ready
$(document).ready(function() {
    $("#input-OnlyBottomBorder").prop('indeterminate', true)
    initServerUrl()
});

// api地址
var SERVER_HOST
function initServerUrl() {
    window.go.gui.App.GetServerUrl().then(async result => {
        SERVER_HOST = result
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}