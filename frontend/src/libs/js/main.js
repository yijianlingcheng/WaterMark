// 隐藏页面元素
function reload() {
    $("#shutterimg-result").html("")
    $("#shutterimg-preview-img").hide()
    $("#shutterimg-preview-img").attr("src", "")
    $("#shutterimg-input").val("")
}

// 显示首页
function showHome() {
    $(".right-content").hide()
    $("#home").show()
}

// 显示快门查询
function showShutterimg() {
    reload()
    $(".right-content").hide()
    $("#shutterimg").show()
}

// 显示水印模板
function showWaterMark() {
    waterMarkBackList()
    $(".right-content").hide()
    $("#watermark").show()
}

// 显示水印生成
function showWaterMarkProcess() {
    $(".right-content").hide()
    $("#watermarkprocess").show()
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

    var data = new FormData();
    data.append('shutterimg', file);
    $.ajax({
        url : "http://localhost:11079/server/getShutterByFile",
        type : "post",
        data : data,
        cache : false,
        processData : false,
        contentType : false,
        success : function (response) {
            var li = "机器快门次数:"+response['MechanicalShutterCount']+" 快门次数:"+response['ShutterCount']
            $("#shutterimg-result").append(li)
        },
        error: function(xhr, status, error) {
            $("#shutterimg-toast").text("查看失败! 原因是:" + error).fadeIn(400).delay(500).fadeOut(400); 
        }
    });
    $("#shutterimg-preview-img").attr("src", "http://localhost:11079/server/getImagePreview?imgagePath="+file+"&random="+Math.random())
    $("#shutterimg-preview-img").show()
    return true
}

$(document).ready(function() {
    $("div.watermark-tpl-item").mouseover(function(){
        $("div.watermark-tpl-item").each(function(){
            $(this).removeClass("shadow p-3 mb-5 bg-white rounded")
        });
        $(this).addClass("shadow p-3 mb-5 bg-white rounded")
    });
});

// 水印模板预览
function waterMarkShowBigImg(object) {
    var src = $(object).attr("src")
    $("#watermark-tpl-detail-img").attr("src", src)
    $(".watermark-tpl-list").hide()
    $("#watermark-tpl-detail").show()
}

// 设置默认水印
function waterMarkSetDefault() {
}

// 返回水印模板列表
function waterMarkBackList() {
    $("#watermark-tpl-detail").hide()
    $(".watermark-tpl-list").show()
    $("#watermark-tpl-detail-img").attr("src", "")
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
            var imgContainer = "<img class='img-imagesContainer' src='http://localhost:11079/server/getImagePreview?imgagePath="+list[0]+"&random="+Math.random() + "'>"
            $("#div-imagesContainer").append(imgContainer)
            $("#div-templateContainer").show()
        }
    }).catch(err => {
        console.log(err);
    }).finally(() => {
        
    });
}

function waterMarkPreivew() {

}

function waterMarkExport() {
    
}