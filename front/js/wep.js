function initPage() {
    if ($.cookie("name") == null) {
        $("#name").text("请登录");
        $("#user").removeAttr("data-toggle");
        $("#user").attr("href", "/login")
    } else {
        $("#name").text($.cookie("name"))
    }
}

function initDashboard() {
    initPage()
    $.base64.utf8decode = true;
    var original = $.base64.atob($("#tData").text());
    if (original == "Empty Set") {
        $("#create").html("<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody><tr><td colspan=\"7\">这里什么都没有</td></tr></tbody>");
        $("#match").html("<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody><tr><td colspan=\"7\">这里什么都没有</td></tr></tbody>");
        return;
    } else {
        var all = original.split("||")
        if (all[0] != "") {
            var cTotal = all[0].split("/");
            var cResult = "<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody>";
            for (let i = 0; i < cTotal.length - 1; i++) {
                const item = cTotal[i].split("|");
                var res = "";
                switch (item[5]) {
                    case "0":
                        res += `<tr class="info">`;
                        break;
                    case "1":
                        res += `<tr class="warning">`;
                        break;
                    case "2":
                        res += `<tr>`;
                        break;
                    case "3":
                        res += `<tr>`;
                        break;
                }
                res += "<td>" + item[0] + "</td>";
                switch (item[1]) {
                    case "0":
                        res += "<td>" + "大作业本" + "</td>";
                        break;
                    case "1":
                        res += "<td>" + "小作业本" + "</td>";
                        break;
                    case "2":
                        res += "<td>" + "英语本" + "</td>";
                        break;
                    case "3":
                        res += "<td>" + "语文作文本" + "</td>";
                        break;
                    case "4":
                        res += "<td>" + "英语作文本" + "</td>";
                        break;
                }
                res += "<td>" + item[2] + "</td>";
                switch (item[3]) {
                    case "0":
                        res += "<td>" + "入手" + "</td>";
                        break;
                    case "1":
                        res += "<td>" + "出手" + "</td>";
                        break;
                }
                res += "<td>" + item[4] + "</td>";
                switch (item[5]) {
                    case "0":
                        res += `<td>待匹配</td><td><button class="btn" type="button" onclick="del('` + item[0] + `')">删除</button></td>`
                        break;
                    case "1":
                        res += `<td>待流通</td><td><button class="btn" type="button" onclick="showInfo('0', '` + item[0] + `')">查看对方信息</button>&nbsp;&nbsp;<button class="btn" type="button" onclick="reject('` + item[0] + `')">拒绝</button>&nbsp;&nbsp;<button class="btn" type="button" onclick="confirm('` + item[0] + `')">确认</button></td>`
                        break;
                    case "2":
                        res += `<td>已完成</td><td><button class="btn" type="button" onclick="showInfo('0', '` + item[0] + `')">查看对方信息</button></td>`
                        break;
                    case "3":
                        res += `<td>已关闭</td><td><button class="btn" type="button" onclick="showInfo('0', '` + item[0] + `')">查看对方信息</button></td>`
                        break;
                }
                cResult += res + "</tr>";
            }
            $("#create").html(cResult + "</tbody>");
        } else {
            $("#create").html("<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody><tr><td colspan=\"7\">这里什么都没有</td></tr></tbody>");
        }
        if (all[1] != "") {
            var mTotal = all[1].split("/");
            var mResult = "<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody>";
            for (let i = 0; i < mTotal.length - 1; i++) {
                const item = mTotal[i].split("|");
                var res = "";
                switch (item[5]) {
                    case "0":
                        res += "<tr class=\"info\">";
                        break;
                    case "1":
                        res += "<tr class=\"warning\">";
                        break;
                    case "2":
                        res += "<tr>";
                        break;
                    case "3":
                        res += "<tr>";
                        break;
                }
                res += "<td>" + item[0] + "</td>";
                switch (item[1]) {
                    case "0":
                        res += "<td>" + "大作业本" + "</td>";
                        break;
                    case "1":
                        res += "<td>" + "小作业本" + "</td>";
                        break;
                    case "2":
                        res += "<td>" + "英语本" + "</td>";
                        break;
                    case "3":
                        res += "<td>" + "语文作文本" + "</td>";
                        break;
                    case "4":
                        res += "<td>" + "英语作文本" + "</td>";
                        break;
                }
                res += "<td>" + item[2] + "</td>";
                switch (item[3]) {
                    case "0":
                        res += "<td>" + "入手" + "</td>";
                        break;
                    case "1":
                        res += "<td>" + "出手" + "</td>";
                        break;
                }
                res += "<td>" + item[4] + "</td>";
                switch (item[5]) {
                    case "0":
                        res += `<td>待匹配</td><td><button class="btn" type="button" onclick="del('` + item[0] + `')">删除</button></td>`
                        break;
                    case "1":
                        res += `<td>待流通</td><td><button class="btn" type="button" onclick="showInfo('1', '` + item[0] + `')">查看对方信息</button>&nbsp;&nbsp;<button class="btn" type="button" onclick="cancel('` + item[0] + `')">取消</button></td>`
                        break;
                    case "2":
                        res += `<td>已完成</td><td><button class="btn" type="button" onclick="showInfo('1', '` + item[0] + `')">查看对方信息</button></td>`
                        break;
                    case "3":
                        res += `<td>已关闭</td><td><button class="btn" type="button" onclick="showInfo('1', '` + item[0] + `')">查看对方信息</button></td>`
                        break;
                }
                mResult += res + "</tr>";
            }
            $("#match").html(mResult + "</tbody>");
        } else {
            $("#match").html("<thead><tr><th>ID</th><th>物品</th><th>数量</th><th>类型</th><th>创建时间</th><th>状态</th><th>操作</th></tr></thead><tbody><tr><td colspan=\"7\">这里什么都没有</td></tr></tbody>");
        }
    }
}

function initMarketPlace() {
    initPage()
    $.base64.utf8decode = true;
    var tables = new Array($("#BigBook tbody"), $("#SmallBook tbody"), $("#EnglishBook tbody"), $("#ChineseComp tbody"), $("#EnglishComp tbody"))
    var original = $.base64.atob($("#tData").text())
    var all = original.split("||")
    for (let i = 0; i < all.length - 1; i++) {
        if (all[i] == "-") {
            tables[i].append(`<tr><td colspan="5" class="text-center">这里什么都没有</td></tr>`)
        } else {
            var items = all[i].split("/")
            for (let j = 0; j < items.length - 1; j++) {
                var item = items[j].split("|")
                if (item[2] == "0") {
                    item[2] = "入手"
                } else {
                    item[2] = "出手"
                }
                var data = `<tr><td>` + item[0] + `</td><td>` + item[1] + `</td><td>` + item[2] + `</td><td>` + item[3] + `</td><td>` + item[4] + `</td><td><button class="btn" type="button" onclick="match('` + item[0] + `')">匹配</button></td></tr>`
                /* alert(data) */
                tables[i].append(data)
            }
        }
    }
}

function create() {
    var json = {
        "uid": $.cookie("uid"),
        "item": $("input[name='item']:checked").val(),
        "amount": $("#amount").val(),
        "kind": $("input[name='kind']:checked").val()
    };
    $.post("/handlers/create", JSON.stringify(json),
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else if (data == "Unauthorized") {
                message = "请先登录"
                $(window).attr('location', '/login');
                return
            } else if (data == "Success") {
                message = "成功创建"
            }
            $("#msg").html(message)
            $('#message').modal()
        });
}

function match(orderID) {
    if ($.cookie("uid") == null) {
        $("#msg").text("请先登录！")
        $('#message').modal()
    } else {
        $.post("/handlers/match", orderID,
            function (data) {
                var message
                if (data == "Server Failure") {
                    message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
                } else if (data == "Unauthorized") {
                    $(window).attr('location', '/login');
                    return
                } else if (data == "Success") {
                    message = "成功匹配"
                } else if (data == "Selfing") {
                    message = "不能匹配自己！"
                } else if (data == "Invalid Status") {
                    message = "看上去已经被抢先匹配了"
                }
                $("#msg").text(message)
                $('#message').modal()
            });
    }
}

function showInfo(role, orderID) {
    var json = {
        "uid": $.cookie("uid"),
        "role": role,
        "orderID": orderID
    };
    $.post("/handlers/oppositeInfo", JSON.stringify(json),
        function (original) {
            var data = JSON.parse(original)
            if (data.status == "Unauthorized") {
                $(window).attr('location', '/login');
                return
            } else if (data.status == "Success") {
                $("#info").html(`姓名：` + data.name + `&nbsp;&nbsp;&nbsp;班级：` + data.class)
            } else if (data.status == "Server Failure") {
                $("#info").html(`服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`)
            }
            $('#oppositeInfo').modal()
        });
}

function del(orderID) {
    $.post("/handlers/delete", orderID,
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else if (data == "Unauthorized") {
                $(window).attr('location', '/login');
                return
            } else if (data == "Success") {
                message = "成功删除"
            }
            $("#msg").text(message)
            $('#message').modal()
        });
}

function reject(orderID) {
    $.post("/handlers/reject", orderID,
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else if (data == "Unauthorized") {
                $(window).attr('location', '/login');
                return
            } else if (data == "Success") {
                message = "成功拒绝"
            }
            $("#msg").text(message)
            $('#message').modal()
        });
}

function cancel(orderID) {
    $.post("/handlers/cancel", orderID,
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else if (data == "Unauthorized") {
                $(window).attr('location', '/login');
                return
            } else if (data == "Success") {
                message = "成功取消"
            }
            $("#msg").text(message)
            $('#message').modal()
        });
}

function confirm(orderID) {
    $.post("/handlers/confirm", orderID,
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else if (data == "Unauthorized") {
                $(window).attr('location', '/login');
                return
            } else if (data == "Success") {
                message = "成功确认"
            }
            $("#msg").text(message)
            $('#message').modal()
        });
}

function exit() {
    $.post("/handlers/exit", $.cookie("uid"),
        function (data) {
            var message
            if (data == "Server Failure") {
                message = `服务器端异常，请<a href="https://github.com/Q5CS/WEP/issues">点此反馈</a>`
            } else {
                $.cookie("name", '', { expires: -1 })
                $.cookie("uid", '', { expires: -1 })
                $.cookie("sessionID", '', { expires: -1 })
                $(window).attr('location', '/');
                return;
            }
            $("#msg").text(message)
            $('#message').modal()
        });
}

function launchFilterModal(index) {
    $("#filterConfirm").click(function () {
        var amount = $("#amount").val();
        var type = $("input[name='type']:checked").val();
        var grade = $("#grade").val();

        var query = "#" + index + " tbody tr"

        $(query).each(function () {
            var a = $(this).children();
            var amo = a[1].innerText;
            var typ = a[2].innerText;
            var gra = a[4].innerText;
            if ((amount != "" && amo != amount) || (grade != "" && gra != grade) || typ != type) {
                a.hide();
            } else {
                a.show();
            }
        })
    })
    $("#filter").modal()
}
