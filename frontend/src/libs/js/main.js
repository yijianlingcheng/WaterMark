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
    $("#shutterimg-preview-img").attr("src", "http://localhost:11079/server/getImagePreview?imgagePath="+file+"&random="+Math.random())
    $("#shutterimg-preview-img").show()
    return true
}

// 获取exif 信息
function getExifInfo(file) {
    var result
    var data = new FormData();
    data.append('shutterimg', file);
    $.ajax({
        url : "http://localhost:11079/server/getShutterByFile",
        type : "post",
        data : data,
        cache : false,
        async : false,
        processData : false,
        contentType : false,
        success : function (response) {
            var li = "机器快门次数:"+response['MechanicalShutterCount']+" 快门次数:"+response['ShutterCount']
            $("#shutterimg-result").append(li)
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
    window.go.gui.App.SelectMultipleImageFile().then(result => {
        if (result.length > 0) {
            $("#div-selectImages").hide()
            $("#watermarkOpenMultipleFiles").val(result)
            var list = result.split(",");
            for (var i = 0; i < list.length; i ++) {
                if (i == 0) {
                    var img = "<img class='img-list img-list-selected pointer' src='http://localhost:11079/server/getImagePreview?imgagePath="+list[i]+"&random="+Math.random() + "'>"
                } else {
                    var img = "<img class='img-list pointer' src='http://localhost:11079/server/getImagePreview?imgagePath="+list[i]+"&random="+Math.random() + "'>"
                }
                $("#watermarkShowMultipleFiles").append(img)
            }
            // 设置模板参数
            setTemplateInfo(getExifInfo(list[0]))
            // 加载模板选项
            loadSelectTemplate()
            // 加载预览图
            loadPreviewImage(list[0], "1", "255,255,255,255", false)
        }
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}

// 加载图片模板类型
function loadSelectTemplate() {
    var data = new FormData();
    $.ajax({
        url : "http://localhost:11079/server/getTplListType",
        type : "get",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            for (var i in response) {
                var li = "<option value="+ i +">"+ response[i] +"</option>"
                $("#select-template").append(li)   
            }
        },
        error: function(xhr, status, error) {
            
        }
    });
}

// 加载实际生成的预览图
function loadPreviewImage(file, tid, color, flag) {
    var data = new FormData();
    data.append("tid", tid)
    data.append("imgagePath", file)
    data.append("color", color)
    data.append("flag", flag.toString())
    $.ajax({
        url : "http://localhost:11079/server/getImageWaterMarkPreview",
        type : "POST",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            // 预览图片
            var imgContainer = "<img class='img-imagesContainer' id='img-imagesContainer' src='http://localhost:11079/server/getImagePreview?imgagePath="+response["SaveImgPath"]+"&random="+Math.random() + "'>"
            // 添加预览图片
            $("#div-imagesContainer").html("").append(imgContainer)
            // 预览容器显示
            $("#div-templateContainer").show()
            // 设置预览参数
            $("#input-Color").val(response["BorderColors"])
            // 保存预览的源文件
            $("#input-PreviewImageFile").val(file)
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
    var File = $("#input-PreviewImageFile").val()
    var Model = $("#input-Model").val()
    var LensModel = $("#input-LensModel").val()
    var Param = $("#input-Params").val()
    var Color = $("#input-Color").val()
    var OnlyBottomFlag = $("#input-OnlyBottomBorder:checked").val() === "on"

    loadPreviewImage(File, "1", Color, OnlyBottomFlag)
}

// 图片导出,js下载存在问题,准备改成Go实现
function waterMarkExport() {
    var imgName = "watermark.jpg"
    var src = $("#img-imagesContainer").attr("src")
    var image = new Image();
    image.src = src;
    image.setAttribute("crossOrigin", "anonymous");
    image.onload = function() {
        let c = document.createElement("canvas");
        c.width = image.width;
        c.height = image.height;
        c.getContext("2d").drawImage(image, 0, 0, image.width, image.height);
        var a = document.createElement("a"); 
        a.download = imgName;
        a.href = c.toDataURL("image/jpg");
        a.click();
    }
}

// ready
$(document).ready(function() {
    $("#input-OnlyBottomBorder").prop('indeterminate', true)
});