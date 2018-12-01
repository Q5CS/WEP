code = request("code");
console.log("code: ", code);

doAuth(code);

function doAuth(code) {
    $.ajax({
        type: "POST",
        url: "/handlers/auth_callback",
        data: {
            "code": code
        },
        dataType: "json",
        success: function (res) {
            console.log(res);
            if(res.status == "fail") {
                alert("错误，请重新登录！");
                $(window).attr('location', '/login');
            }
            $.cookie("name", res.name)
            $.cookie("uid", res.uID)
            $.cookie("sessionID", res.sessionID)
            $(window).attr('location', '/dashboard');
        }
    });
}

function request(paras) {
    var url = location.href;
    var paraString = url.substring(url.indexOf("?") + 1, url.length).split("&");
    var paraObj = {}
    for (i = 0; j = paraString[i]; i++) {
        paraObj[j.substring(0, j.indexOf("=")).toLowerCase()] = j.substring(j.indexOf("=") + 1, j.length);
    }
    var returnValue = paraObj[paras.toLowerCase()];
    if (typeof (returnValue) == "undefined") {
        return "";
    } else {
        return returnValue;
    }
}