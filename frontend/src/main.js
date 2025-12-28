// 获取服务端接口地址
var host = "";
function setServerHost(h) {
    host = h
}

(function () {
    window.go.ui.App.GetApiServerHost().then(result => {
        setServerHost(result)
    }).catch(err => {
        console.log(err);
    }).finally(() => {

    });
})()

// 休眠
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
}

// 获取年月日时分秒
function getDate() {
    let currentDate = new Date();

    // 获取年、月、日、小时、分钟和秒
    let year = currentDate.getFullYear();
    let month = (currentDate.getMonth() + 1).toString().padStart(2, '0'); // 月份从0开始，所以要加1，并补零
    let day = currentDate.getDate().toString().padStart(2, '0'); // 补零
    let hours = currentDate.getHours().toString().padStart(2, '0'); // 补零
    let minutes = currentDate.getMinutes().toString().padStart(2, '0'); // 补零
    let seconds = currentDate.getSeconds().toString().padStart(2, '0'); // 补零

    // 拼接日期和时间字符串
    return year + '-' + month + '-' + day + ' ' + hours + ':' + minutes + ':' + seconds;
}

/*****************开始*******************/
// 下面是启动页方法

// 加载启动页
async function loadHomePageMessage() {
    while (host == "") {
        await sleep(50)
    }
    $.get(host + "/message/info", async function (data) {
        let jsonData = JSON.parse(data)
        for (let i = 0; i < jsonData.list.length; i++) {
            await sleep(10)
            $("#g-start-info").html(jsonData.list[i] + "...")
        }
        if (jsonData.list[jsonData.list.length - 1] == "启动完成") {
            await sleep(10)
            window.location.href = "/A/aboutVersionView.html"
        }
    })
}
/*****************结束*******************/


/*****************开始*******************/
// 下面是exif查看页面方法

// exif 信息展示操作
const ExifInfoProcess = {
    // 清空展示的exif信息
    exifCleanExifInfo: function (file) {
        $("#exifinfo-export").hide()
        $("#exifinfo-export-input").val("")
        $("#exifinfo-file").html(file)
    },
    // 将exif信息置为获取失败
    exifFailureExifInfo: function () {
        $("#exifinfo-e").html("获取失败")
        $("#exifinfo-m").html("获取失败")
        $("#exifinfo-p").html("获取失败")
        $("#exifinfo-f").html("获取失败")
        $("#exifinfo-t").html("获取失败")
        $("#exifinfo-s").html("获取失败")
    },
    // 设置并展示对应的exif信息
    exifSetExifInfoAndShow: function (exifInfo) {
        $("#exifinfo-e").html(exifInfo["器材"])
        $("#exifinfo-m").html(exifInfo["模式"])
        $("#exifinfo-p").html(exifInfo["参数"])
        $("#exifinfo-f").html(exifInfo["焦距"])
        $("#exifinfo-t").html(exifInfo["时间"])
        $("#exifinfo-s").html(exifInfo["快门次数"])

        // 判断照片角度,防止竖构图照片无法正常显示
        if (exifInfo["Orientation"].indexOf("Horizontal") > -1 || exifInfo["Orientation"].indexOf("normal") > -1) {
            // 兼容处理
            if (parseInt(exifInfo["ImageHeight"]) > parseInt(exifInfo["ImageWidth"])) {
                $("#exifinfo-preview-img").attr("src", "").removeClass().addClass("exifinfo-preview-img-height")
            } else {
                $("#exifinfo-preview-img").attr("src", "").removeClass().addClass("exifinfo-preview-img-width")
            }
        } else {
            $("#exifinfo-preview-img").attr("src", "").removeClass().addClass("exifinfo-preview-img-height")
        }
        $("#exifinfo-preview-img").attr("src", host + "/view/showImage?file=" + exifInfo["file"])
        $("#exifinfo-export-input").val(exifInfo["file"])
        $("#exifinfo-export").show()
    },
    // 选择图片
    exifViewSelectImage: function () {
        window.go.ui.App.SelectImageFile().then(result => {
            ExifImageOptions.getPhotoExifInfo(result)
        }).catch(err => {
            console.log(err);
        }).finally(() => {

        });
    },
}

// exif查看页面对应的操作
const ExifImageOptions = {
    // 获取照片exif信息
    getPhotoExifInfo: function (path) {
        let data = new FormData();
        data.append('file', path);
        $.ajax({
            url: host + "/view/getImagesExifInfo",
            type: "post",
            data: data,
            cache: false,
            processData: false,
            contentType: false,
            success: function (response) {
                ExifInfoProcess.exifCleanExifInfo(response["file"])
                // 获取exif失败
                if (response["code"] > 0) {
                    ExifInfoProcess.exifFailureExifInfo()
                    return
                }
                ExifInfoProcess.exifSetExifInfoAndShow(response)
            },
            error: function (xhr, status, error) {
                //TODO 提示错误
            }
        });
    },
    // 导出照片exif信息,使用后端生成文件方式
    exportPhotoExifInfo: function () {
        let file = $("#exifinfo-export-input").val()
        if (file.length > 0) {
            window.go.ui.App.ShowSaveFileDialog("exifInfo.csv").then(result => {
                let data = new FormData();
                data.append('file', file);
                data.append('save', result);
                fetch(host + "/view/exifInfoExportv2", {
                    method: 'POST',
                    body: data
                })
                    .then(response => {
                        response.json().then(data => {
                            console.log(data)
                        })
                    })
            }).catch(err => {
                console.log(err);
            }).finally(() => {

            });
        }
    }
}


/*****************结束*******************/

/*****************开始*******************/

var defaultBucket = ""
var PreviewBucket = "preview"

const FrameViewGoEvent = {
    // 清理暂存的信息
    TemporaryClean: async function () {
        window.go.ui.App.TemporaryClean().then(result => {
            console.log(result)
        }).catch(err => {
            console.log(err);
        }).finally(() => {

        });
    },
    // 将指定的key,value暂存到指定桶中
    TemporaryStorage: async function (key, value, bucket) {
        window.go.ui.App.TemporaryStorage(key, value, bucket).then(result => {

        }).catch(err => {
            console.log(err);
        }).finally(() => {

        });
    },
    // 根据指定key与桶获取暂存的内容
    TemporaryGet: async function (key, bucket) {
        return window.go.ui.App.GetTemporaryStorage(key, bucket).then(result => {
            return result
        }).catch(err => {
            console.log(err)
            return ""
        }).finally(() => {

        });
    },
    // 获取指定桶暂存的全部内容
    GetTemporaryAll: async function (bucket) {
        return window.go.ui.App.GetTemporaryAll(bucket).then(result => {
            return result
        }).catch(err => {
            console.log(err)
            return ""
        }).finally(() => {

        });
    },
    // 展示导出提示
    SelectDirectory: async function (title) {
        return window.go.ui.App.SelectDirectory(title).then(result => {
            return result
        }).catch(err => {
            console.log(err)
            return ""
        }).finally(() => {

        });
    },
    // 获取文件夹中的jpg图片
    GetDirectoryJpgFiles: async function name(path) {
        return window.go.ui.App.GetDirectoryJpgFiles(path).then(result => {
            return result
        }).catch(err => {
            console.log(err)
            return ""
        }).finally(() => {

        });
    },
    // 展示导出确认提示
    SureExportPhotoTips: async function (title) {
        return window.go.ui.App.ShowExportPhotoTips(title).then(result => {
            return result
        }).catch(err => {
            console.log(err)
            return ""
        }).finally(() => {

        });
    },
}

const FrameViewSelectImageOptions = {
    // 选择图片文件,多选
    SelectImages: function () {
        window.go.ui.App.SelectMultipleImageFile().then(result => {
            FrameViewGoEvent.TemporaryStorage("frame-select-images", result, defaultBucket)
            FrameViewSelectImageProcess.AddAndShowSelectImages(result)
        }).catch(err => {
            console.log(err);
        }).finally(() => {

        });
    },
    // 选择图片文件夹
    SelectDirImages: async function () {
        let dir = await FrameViewGoEvent.SelectDirectory("请选择文件夹")
        if (dir != "") {
            let file = await FrameViewGoEvent.GetDirectoryJpgFiles(dir)
            if (file != "") {
                FrameViewGoEvent.TemporaryStorage("frame-select-images", file, defaultBucket)
                FrameViewSelectImageProcess.AddAndShowSelectImages(file)
            }
        }
    },
    // 跳转至选择模板
    LocationSelectTpl: function () {
        window.location.href = "/F/selectTpl.html"
    },
}

const FrameViewSelectImageProcess = {
    // 填充展示选中的文件名称
    AddAndShowSelectImages: function (str) {
        if (str == "") {
            return
        }
        let list = str.split(",")
        $("#select-file-num").html("共选择了" + list.length + "个文件")
        $("#file-list-show-ul").html("")
        let li = ""

        for (let i = 0; i < list.length; i++) {
            if (i > 25) {
                li += "<li> ...... </li>"
                break
            }
            li += "<li> " + list[i] + " </li>"
        }
        $("#file-list-show-ul").html(li)
        $(".layout-tpl").css("display", "block")
    },
}

const FrameViewSelectTplOptions = {
    // 加载模板信息
    LoadTemplates: async function () {
        while (host == "") {
            await sleep(50)
        }
        $.get(host + "/frame/getFrameTemplateInfo", function (data) {
            if (data["code"] == 0) {
                FrameViewSelectTplProcess.AddAndShowTplImages(data["list"])
            }
        })
    },
    // 确认是否选择模板,并执行跳转方法
    LocationConfirmTpl: async function () {
        len = $(".tpl-selected").length
        if (len == 0) {
            $("div.active").addClass("shake")
            await sleep(1000)
            $("div.active").removeClass("shake")
            return
        }
        let tplName = $(".tpl-selected").attr("data-src")
        let tplLink = $(".tpl-selected>img").attr("src")

        FrameViewGoEvent.TemporaryStorage("frame-select-tpl-name", tplName, defaultBucket)
        FrameViewGoEvent.TemporaryStorage("frame-select-tpl-image", tplLink, defaultBucket)
        window.location.href = "/F/confirmTpl.html"
    },
    // 选中模板
    AddClassOnSelectTpl: function (obj) {
        $(".slide").removeClass("tpl-selected")
        $(obj).addClass("tpl-selected")
    },
}

const FrameViewSelectTplProcess = {
    // 填充展示模板信息
    AddAndShowTplImages: function (data) {
        let index = 0
        let names = ""
        let images = ""
        for (let key in data) {
            index++
            if (index == 1) {
                $("#tpl-show-name-span").html(key)

                names += "<div class=\"dot active\"><span>" + key + "</span></div>"
                images += "<div data-src=\"" + key + "\" onclick=\"FrameViewSelectTplOptions.AddClassOnSelectTpl(this)\" class=\"slide active\" style=\"left: 0px;\"><img class=\"images\" src=\"" + host + "/view/showImage?file=" + data[key] + "\"></div>"

                continue
            }
            names += "<div class=\"dot\"><span>" + key + "</span></div>"
            images += "<div data-src=\"" + key + "\" onclick=\"FrameViewSelectTplOptions.AddClassOnSelectTpl(this)\" class=\"slide \" style=\"left: 0px;\"><img class=\"images\" src=\"" + host + "/view/showImage?file=" + data[key] + "\"></div>"
        }
        $("#tpl-show-name-div").html(names)
        $("#tpl-show-images-div").html(images)

        const track = document.querySelector('.track');
        const slides = Array.from(track.children);
        const prevBtn = document.querySelector('.btn.btn-back');
        const nextBtn = document.querySelector('.btn.btn-next');
        const navIndicator = document.querySelector('.nav-indicator');
        const dots = Array.from(navIndicator.children)
        const slideSize = slides[0].getBoundingClientRect();
        const slideWidth = slideSize.width;

        var tl = new TimelineMax();
        function blur(el, blur) {
            tl.fromTo(el, 0.55,
                { filter: `blur(${blur}px)` },
                { filter: 'blur(0px)' });
        }

        const slidePosition = (slide, index) => {
            slide.style.left = `${slideWidth * index}px`;
        }
        slides.forEach(slidePosition)

        const slideToMove = (track, currentSlide, targetSlide) => {
            track.style.transform = `translateX(-${targetSlide.style.left})`;
            currentSlide.classList.remove('active');
            targetSlide.classList.add('active');
            $("#tpl-show-name-span").html($(targetSlide).attr("data-src"))
        }

        function updateDots(current, target) {
            current.classList.remove('active')
            target.classList.add('active')
        }

        function btnShowHide(targetIndex, prevBtn, nextBtn, slides) {
            if (targetIndex == 0) {
                prevBtn.classList.add('hidden')
                nextBtn.classList.remove('hidden')
            } else if (targetIndex == slides.length - 1) {
                prevBtn.classList.remove('hidden')
                nextBtn.classList.add('hidden')
            } else {
                prevBtn.classList.remove('hidden')
                nextBtn.classList.remove('hidden')
            }
        }

        nextBtn.addEventListener('click', (e) => {
            var currentSlide = track.querySelector('.active')
            var nextSlide = currentSlide.nextElementSibling;
            var currentDot = navIndicator.querySelector('.active');
            var nextDot = currentDot.nextElementSibling;
            var nextIndex = slides.findIndex(slide => slide === nextSlide)

            slideToMove(track, currentSlide, nextSlide);
            updateDots(currentDot, nextDot);
            btnShowHide(nextIndex, prevBtn, nextBtn, slides);
            if (e.detail > 1) return;
            blur(track, 5)
        });

        prevBtn.addEventListener('click', (e) => {
            var currentSlide = track.querySelector('.active')
            var prevSlide = currentSlide.previousElementSibling;
            var currentDot = navIndicator.querySelector('.active');
            var prevDot = currentDot.previousElementSibling;
            var prevIndex = slides.findIndex(slide => slide === prevSlide)

            slideToMove(track, currentSlide, prevSlide);
            updateDots(currentDot, prevDot);
            btnShowHide(prevIndex, prevBtn, nextBtn, slides)
            if (e.detail > 1) return;
            blur(track, 5)
        });

        navIndicator.addEventListener('click', (e) => {
            var targetDot = e.target.closest('.dot');
            if (!targetDot) return;

            var currentSlide = track.querySelector('.active');
            var currentDot = navIndicator.querySelector('.active');
            var targetIndex = dots.findIndex(dot => dot === targetDot)
            var targetSlide = slides[targetIndex];

            slideToMove(track, currentSlide, targetSlide)
            updateDots(currentDot, targetDot);
            btnShowHide(targetIndex, prevBtn, nextBtn, slides)
            if (e.detail > 1) return;
            blur(track, 5)
        });
    },
}

const FrameViewConfirmTplOptions = {
    LoadConfirmTpl: async function () {
        let name = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)
        let image = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-image", defaultBucket)
        let images = "<div data-src=\"" + name + "\" onclick=\"FrameViewSelectTplOptions.AddClassOnSelectTpl(this)\" class=\"slide active\" style=\"left: 0px;\"><img class=\"images\" src=\"" + image + "\"></div>"

        $("#tpl-show-name-span").html("已选择布局:" + name)
        $("#tpl-show-images-div").html(images)
    },
    // 跳转展示边框预览页面
    SubmitAndPreviewBorder: async function () {
        window.location.href = "/F/previewBorder.html"
    },
    // 导出提示
    ExportPhotoTips: async function (from) {
        let exportPath = await FrameViewGoEvent.SelectDirectory("选择存放导出图片的目录")
        if (exportPath == "") {
            return
        }
        await sleep(100)
        let sure = await FrameViewGoEvent.SureExportPhotoTips(exportPath)
        if (sure == "Yes" || sure == "是") {
            console.log(sure, exportPath)
            FrameViewExportProcess.LocationExport(exportPath)
        }
    }
}


var framePreviewSelectImages = []
var lastOptime = Date.now();
var minOpTime = 400
var inUpdate = false
var requestTime = Date.now()
var requestUseTime = 0

// 是否是快速点击
function isfastClick() {
    let currentTime = Date.now();
    if (currentTime - lastOptime < minOpTime) {
        lastOptime = currentTime
        return true
    }
    lastOptime = currentTime
    return false
}

// 只检查
function onlyCheckisfastClick() {
    let currentTime = Date.now();
    return (currentTime - lastOptime) < minOpTime
}

const FramePreviewBorderOptions = {
    // 设置选中的可预览图片
    SetSelectImages: async function () {
        let files = await FrameViewGoEvent.TemporaryGet("frame-select-images", defaultBucket)
        framePreviewSelectImages = files.split(",")
    },
    // 加载预览图片
    LoadPreivewImage: async function (name, layout) {
        // 记录接口请求时间
        requestTime = Date.now()

        inUpdate = true
        while (host == "") {
            await sleep(50)
        }
        let files = await FrameViewGoEvent.TemporaryGet("frame-select-images", defaultBucket)
        let fileList = files.split(",")

        if (name == "") {
            name = fileList[0]
        }

        let preivewStr = await FrameViewGoEvent.TemporaryGet(name, PreviewBucket)
        if (preivewStr != "") {
            layout = preivewStr
        } else {
            if (layout == "") {
                let layoutName = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)
                layout = '{"frame_name":"' + layoutName + '"}'
            }
        }

        let data = new FormData();
        data.append('file', name);
        data.append('type', "border");
        data.append('layout', layout)

        // 发起 POST 请求
        fetch(host + "/frame/showPhotoFrame", {
            method: 'POST',
            body: data
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                // 将响应体解析为 Blob
                return response.blob();
            })
            .then(blob => {
                // 记录请求耗时
                requestUseTime = Date.now() - requestTime

                inUpdate = false
                $("#images-container").attr("data-src", name)
                // 创建指向 Blob 的临时 URL
                const imageUrl = URL.createObjectURL(blob);
                // 创建 <img> 元素并设置 src
                const img = document.createElement('img');
                document.getElementById("images-container").src = imageUrl;
                // 可选：在图片加载完成后释放 Blob URL 内存
                img.onload = () => {
                    URL.revokeObjectURL(imageUrl); // 防止内存泄漏
                };

                $("#images-origin").attr("src", host + "/view/showImage?file=" + name)
                $("#blur-background-img").attr("src", host + "/view/showImage?file=" + name)
            })
            .catch(error => {
                console.error('获取或显示图片失败:', error);
            });
    },
    // 加载展示的exif信息
    LoadPreivewExif: async function (name) {
        if (name == "") {
            let files = await FrameViewGoEvent.TemporaryGet("frame-select-images", defaultBucket)
            name = files.split(",")[0]
        }
        let layoutName = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)

        let data = new FormData();
        data.append('file', name);
        data.append('type', "border");
        data.append('layout', '{"frame_name":"' + layoutName + '"}')

        fetch(host + "/frame/getExifAndBorderInfo", {
            method: 'POST',
            body: data
        })
            .then(response => {
                response.json().then(data => {
                    if (data["code"] == 0) {
                        FramePreviewBorderDomProcess.UpdateExifInfo(data["exif"])
                        FramePreviewBorderDomProcess.UpdateTextInfo(data["text"])
                        FramePreviewBorderDomProcess.UpdateOriginImageCss(data["size"])
                    }
                })
            })
    },
}

var lastWidthAndHeightFlag = -1
// 是否首次加载
function isFirstLoad() {
    return lastWidthAndHeightFlag == -1
}
// 是否宽大于高
function isWidthGtHeight() {
    return lastWidthAndHeightFlag > 0
}
// 设置宽大于高
function setWithGtHeight() {
    lastWidthAndHeightFlag = 1
}
// 设置宽小于高
function setWidthLtHeight() {
    lastWidthAndHeightFlag = 0
}


const FramePreviewBorderDomProcess = {
    // 加载原始图片
    LoadOriginImage: async function () {
        while (host == "") {
            await sleep(50)
        }
        let files = await FrameViewGoEvent.TemporaryGet("frame-select-images", defaultBucket)
        let list = files.split(",")
        let images = ""
        for (let i = 0; i < list.length; i++) {
            let tmp = list[i].split("/")
            let name = tmp[tmp.length - 1]
            if (i <= 5) {
                if (i == 0) {
                    images += '<div class="div-img-list"><img class="img-list img-list-selected" data-src="' + list[i] + '" src="' + host + "/view/showImage?file=" + list[i] + '" onclick="FramePreviewBorderDomProcess.UpdatePrevireImage(this)"><span class="img-list-span">' + name + '</span></div>'
                    continue
                }
                images += '<div class="div-img-list"><img class="img-list" data-src="' + list[i] + '" src="' + host + "/view/showImage?file=" + list[i] + '" onclick="FramePreviewBorderDomProcess.UpdatePrevireImage(this)"><span class="img-list-span">' + name + '</span></div>'
            }
        }
        $("#origin-images-container").html(images)
    },
    // 更新展示的exif信息
    UpdateExifInfo: function (data) {
        $("#exif-info-name").html(data["FileName"])
        $("#exif-info-byte").html((parseInt(data["ImageDataSize"]) / 1024 / 1024).toFixed(2) + "MB")
        $("#exif-info-size").html(data["ImageSize"])
        $("#exif-info-model").html(data["Model"])
        $("#exif-info-len").html(data["LensModel"])
        $("#exif-info-focal").html(data["FocalLength"])
        $("#exif-info-fnumber").html(data["FNumber"])
        $("#exif-info-exposure").html(data["ExposureTime"])
        $("#exif-info-iso").html(data["ISO"])
        $("#exif-info-time").html(data["时间"])
    },
    // 保存预览时修改的参数
    SaveImageParams: async function () {
        let newData = []
        let counter = 0
        let len = $("input[type=text]").length
        $("input[type=text]").each(async function () {
            counter++
            let name = $(this).attr("name")
            let key = $(this).attr("data-name")
            let value = $(this).val()
            newData.push(name)
            newData.push(key)
            newData.push(value)
            if (counter >= len) {
                let layoutName = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)
                let imageName = $("#images-container").attr("data-src")
                let layout = {
                    frame_name: layoutName,
                }
                for (let i = 0; i < newData.length; i = i + 3) {
                    layout[newData[i]] = newData[i + 2]
                }
                FrameViewGoEvent.TemporaryStorage(imageName, JSON.stringify(layout), PreviewBucket)
            }
        })
    },
    // 更新展示的参数
    UpdateTextInfo: async function (data) {
        let imageName = $(".img-list-selected").attr("data-src")
        while (imageName == "" || imageName == undefined) {
            await sleep(10)
            imageName = $(".img-list-selected").attr("data-src")
        }
        let preivewStr = await FrameViewGoEvent.TemporaryGet(imageName, PreviewBucket)
        let preivewParams = {}
        if (preivewStr != "") {
            preivewParams = JSON.parse(preivewStr)
        }
        let html = ""
        for (let i = 0; i < data.length; i = i + 3) {
            let name = data[i]
            let key = data[i + 1]
            let value = data[i + 2]
            if (name in preivewParams) {
                value = preivewParams[name]
            }
            html += "<div><p>" + key + "</p><input type='text' spellcheck='false' name='" + name + "' data-name='" + key + "' value='" + value + "' /></div>"
        }
        let saveInfo = "<button id='save-btn' onclick='FramePreviewBorderDomProcess.SaveImageParams()'>保存参数</button>"
        let exportFile = "<button id='export-btn' onclick='FrameViewConfirmTplOptions.ExportPhotoTips(\"preview\")'>导出全部</button>"
        html += "<div class='show-text-options'>" + saveInfo + exportFile + "<div>"
        $("#show-text-div").html(html)
        $("input[type=text]").blur(function () {
            let newData = []
            let counter = 0
            let len = $("input[type=text]").length
            $("input[type=text]").each(async function () {
                counter++
                let name = $(this).attr("name")
                let key = $(this).attr("data-name")
                let value = $(this).val()
                newData.push(name)
                newData.push(key)
                newData.push(value)
                if (counter >= len) {
                    let isSame = JSON.stringify(data) === JSON.stringify(newData)
                    if (!isSame) {
                        // 重新赋值
                        data = newData
                        let layoutName = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)
                        let layout = {
                            frame_name: layoutName,
                        }
                        for (let i = 0; i < newData.length; i = i + 3) {
                            layout[newData[i]] = newData[i + 2]
                        }
                        FramePreviewBorderOptions.LoadPreivewImage(imageName, JSON.stringify(layout))
                    }
                }
            })
        });
    },
    // 更新预览中的原始图片的宽高
    UpdateOriginImageCss: async function (data) {
        let showWidth = 500
        let showHeight = 381

        let borderLeft = parseInt(data["BorderLeftWidth"])
        let borderRight = parseInt(data["BorderRightWidth"])
        let borderTop = parseInt(data["BorderTopHeight"])
        let borderBottom = parseInt(data["BorderBottomHeight"])
        let width = parseInt(data["SourceWidth"])
        let height = parseInt(data["SourceHeight"])
        let borderRadius = parseInt(data["BorderRadius"])
        let totalWidth = borderLeft + borderRight + width
        let totalHeight = borderTop + borderBottom + height

        if (width >= height) {

            let marginLeft = (borderLeft / totalWidth * showWidth).toFixed()
            let showHeight = (showWidth / totalWidth * totalHeight).toFixed()
            let marginRight = (borderRight / totalWidth * showWidth).toFixed()
            let marginTop = (borderTop / totalHeight * showHeight).toFixed()

            let sleepTime = 0
            if (!isFirstLoad() && !isWidthGtHeight()) {
                sleepTime = parseInt(requestUseTime * 2)

            }
            setWithGtHeight()
            if (sleepTime > 0) {
                $("#images-container").hide()
                $("#images-origin").hide()
                $("#blur-show-div").hide()

                await sleep(parseInt(requestUseTime * 2))
            }
            // 原始图片展示宽度
            let originShowWidth = showWidth - marginLeft - marginRight
            let originCss = { "margin-top": marginTop + "px", "margin-left": (marginLeft - showWidth) + "px", "width": (originShowWidth) + "px", "height": "auto" }
            if (borderRadius > 0) {
                let radius = parseInt(originShowWidth / width * borderRadius)
                originCss["border-radius"] = radius + "px"
                originCss["box-shadow"] = "0px 0px " + radius + "px " + radius / 2 + "px rgba(128, 128, 128, 0.5)"
            }
            $("#images-container").css({ "height": "auto", "width": showWidth + "px" })
            $("#images-origin").css(originCss)

            if (borderRadius > 0) {
                $("#blur-show-div").css({ "height": showHeight + "px", "width": showWidth + "px", "margin-left": 0 + "px" })
            }

            if (sleepTime > 0) {
                $("#images-container").show()
                $("#images-origin").show()
            }

            if (borderRadius > 0) {
                $("#blur-show-div").show()
            }

        } else {

            let newShowWidth = showHeight / totalHeight * totalWidth
            let marginLeft = borderLeft / totalWidth * newShowWidth - newShowWidth

            let newOriginShowHegiht = showHeight / totalHeight * height
            let newOriginShowWidth = showHeight / totalHeight * width

            let marginTop = (borderTop / totalHeight * showHeight).toFixed()
            let sleepTime = 0
            if (!isFirstLoad() && isWidthGtHeight()) {
                sleepTime = parseInt(requestUseTime * 2)
            }
            setWidthLtHeight()

            if (sleepTime > 0) {
                $("#images-container").hide()
                $("#images-origin").hide()
                $("#blur-show-div").hide()

                await sleep(parseInt(requestUseTime * 2))
            }

            let originCss = { "height": newOriginShowHegiht + "px", "width": newOriginShowWidth + "px", "margin-left": marginLeft + "px", "margin-top": marginTop + "px" }
            if (borderRadius > 0) {
                let radius = parseInt(newOriginShowHegiht / height * borderRadius)
                originCss["border-radius"] = radius + "px"
                originCss["box-shadow"] = "0px 0px " + radius + "px " + radius / 2 + "px rgba(128, 128, 128, 0.5)"
            }
            let containerCss = { "height": showHeight + "px", "width": newShowWidth + "px" }
            // 调整模糊模板的容器大小
            $("#images-container").css(containerCss)
            if (borderRadius > 0) {
                containerCss["margin-left"] = (showWidth - newShowWidth) / 2 + "px"
                $("#blur-show-div").css(containerCss)
            }

            $("#images-origin").css(originCss)

            if (sleepTime > 0) {
                $("#images-container").show()
                $("#images-origin").show()
            }
            if (borderRadius > 0) {
                $("#blur-show-div").show()
            }
        }
    },
    // 更新选中的预览图片
    UpdatePrevireImage: async function (object) {
        $(".img-list").removeClass("img-list-selected")
        $(object).addClass("img-list-selected")
        let name = $(object).attr("data-src")

        FramePreviewBorderOptions.LoadPreivewExif(name)
        FramePreviewBorderOptions.LoadPreivewImage(name, "")
    },
    // 选中前一个元素
    PrevPreviewImage: async function () {
        // 获取当前选中的元素
        let object = $(".img-list-selected")[0]
        let imageLink = $(".img-list-selected").attr("data-src")
        const parent = object.parentNode.parentNode;
        const childrenArray = Array.from(parent.children);
        const index = childrenArray.indexOf(object.parentNode)
        let indexOfArr = framePreviewSelectImages.indexOf(imageLink)
        if (index == 0 && indexOfArr == 0) {
            return
        }
        if (index > 0) {

            let preIndex = index - 1
            $(".img-list").removeClass("img-list-selected")
            $(".img-list:eq(" + preIndex + ")").addClass("img-list-selected")

            if (isfastClick()) {
                return
            }
            let name = $(".img-list-selected").attr("data-src")
            FramePreviewBorderOptions.LoadPreivewExif(name)
            FramePreviewBorderOptions.LoadPreivewImage(name, "")
            return
        }
        if (index == 0) {

            let preIndex = indexOfArr - 1
            $(".img-list").removeClass("img-list-selected")
            let imageIndex = 0
            for (let i = preIndex; i <= preIndex + 5; i++) {
                let imageNameArr = framePreviewSelectImages[i].split("/")
                let imageName = imageNameArr[imageNameArr.length - 1]
                $(".img-list:eq(" + imageIndex + ")").attr("data-src", framePreviewSelectImages[i])
                $(".img-list-span:eq(" + imageIndex + ")").html(imageName)
                $(".img-list:eq(" + imageIndex + ")").attr("src", host + "/view/showImage?file=" + framePreviewSelectImages[i])
                imageIndex++
            }
            $(".img-list:eq(0)").addClass("img-list-selected")
            if (isfastClick()) {
                return
            }
            let name = $(".img-list-selected").attr("data-src")
            FramePreviewBorderOptions.LoadPreivewExif(name)
            FramePreviewBorderOptions.LoadPreivewImage(name, "")
            return
        }
    },
    // 选中后一个元素
    NextPreviewImage: async function () {
        // 获取当前选中的元素
        let object = $(".img-list-selected")[0]
        let imageLink = $(".img-list-selected").attr("data-src")

        const parent = object.parentNode.parentNode;
        const childrenArray = Array.from(parent.children);
        const index = childrenArray.indexOf(object.parentNode)

        let indexOfArr = framePreviewSelectImages.indexOf(imageLink)

        if (index == 5 && indexOfArr == framePreviewSelectImages.length - 1) {
            return
        }
        if (index == framePreviewSelectImages.length - 1) {
            return
        }
        if (index < 5) {

            let preIndex = index + 1
            $(".img-list").removeClass("img-list-selected")
            $(".img-list:eq(" + preIndex + ")").addClass("img-list-selected")

            if (isfastClick()) {
                return
            }

            let name = $(".img-list-selected").attr("data-src")

            FramePreviewBorderOptions.LoadPreivewExif(name)
            FramePreviewBorderOptions.LoadPreivewImage(name, "")

            return
        }
        if (index == 5) {

            let preIndex = indexOfArr + 1
            $(".img-list").removeClass("img-list-selected")
            let imageIndex = 0
            for (let i = preIndex - 5; i <= preIndex; i++) {
                let imageNameArr = framePreviewSelectImages[i].split("/")
                let imageName = imageNameArr[imageNameArr.length - 1]
                $(".img-list:eq(" + imageIndex + ")").attr("data-src", framePreviewSelectImages[i])
                $(".img-list-span:eq(" + imageIndex + ")").html(imageName)
                $(".img-list:eq(" + imageIndex + ")").attr("src", host + "/view/showImage?file=" + framePreviewSelectImages[i])
                imageIndex++
            }
            $(".img-list:eq(5)").addClass("img-list-selected")

            if (isfastClick()) {
                return
            }

            let name = $(".img-list-selected").attr("data-src")

            FramePreviewBorderOptions.LoadPreivewExif(name)
            FramePreviewBorderOptions.LoadPreivewImage(name, "")

            return
        }
    },
    // 检查当前选中的图片是否与预览展示的一直
    CheckSelectAndPreviewIsSame: async function () {
        if (inUpdate) {
            return
        }
        let select = $(".img-list-selected").attr("data-src")
        let preivew = $("#images-container").attr("data-src")

        if (select == preivew && !onlyCheckisfastClick()) {
            return
        }
        inUpdate = true

        // 加载预览大图
        FramePreviewBorderOptions.LoadPreivewImage(select, "")
        // 加载大图的exif信息
        FramePreviewBorderOptions.LoadPreivewExif(select)
        // 更新
        isfastClick()

        inUpdate = false
    },
}

/***********结束************/

/***********开始************/
// 导出
const FrameViewExport = {
    // 跳转至导出页面
    LocationToExport: async function (object) {
        let savePath = $(object).attr("data-src")
        FrameViewExportProcess.LocationExport(savePath)
    },
    // 导出
    ExportPhoto: async function () {
        while (host == "") {
            await sleep(50)
        }
        let save = await FrameViewGoEvent.TemporaryGet("export-save-path", defaultBucket)
        let files = await FrameViewGoEvent.TemporaryGet("frame-select-images", defaultBucket)
        let layoutName = await FrameViewGoEvent.TemporaryGet("frame-select-tpl-name", defaultBucket)
        if (save == "" || files == "") {
            return
        }
        let preivewParams = await FrameViewGoEvent.GetTemporaryAll(PreviewBucket)
        layout = '{"frame_name":"' + layoutName + '"}'

        let data = new FormData();
        data.append('save', save);
        data.append('file', files);
        data.append('layout', layout);
        data.append('preview_layout', preivewParams)
        fetch(host + "/frame/createExportTask", {
            method: 'POST',
            body: data
        })
            .then(response => {
                response.json().then(data => {
                    FrameViewGoEvent.TemporaryClean()
                    FrameViewExportProcess.ExportProgress(files.split(",").length)
                })
            })
    },
}

const FrameViewExportProcess = {
    // 跳转
    LocationExport: async function (savePath) {
        FrameViewGoEvent.TemporaryStorage("export-save-path", savePath, defaultBucket)
        window.location.href = "/F/exportView.html"
    },
    // 导出进度条
    ExportProgress: async function (fileNum) {
        let currentNum = 0
        while (host == "") {
            await sleep(50)
        }
        const eventSource = new EventSource(host + '/frame/getExportProgress');

        eventSource.onmessage = (event) => {
            let files = event.data.split(",")
            currentNum += files.length
            let ratio = parseInt(currentNum / fileNum * 100)
            console.log(ratio, currentNum, fileNum)
            $("#export-progress-label").html("导出进度:" + ratio + "%")
            $("#export-progress").attr("value", ratio)
            if (fileNum <= currentNum) {
                eventSource.close()
            }
        };
    }
}