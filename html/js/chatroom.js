function setCookie(c_name, value, expiredays) {
  var exdate = new Date()
  exdate.setDate(exdate.getDate() + expiredays)
  document.cookie = c_name + "=" + escape(value) +
    ((expiredays == null) ? "" : ";expires=" + exdate.toGMTString()) + "; path=/"
}

function getCookie(c_name) {
  if (document.cookie.length > 0) {
    c_start = document.cookie.indexOf(c_name + "=")
    if (c_start != -1) {
      c_start = c_start + c_name.length + 1
      c_end = document.cookie.indexOf(";", c_start)
      if (c_end == -1) c_end = document.cookie.length
      return unescape(document.cookie.substring(c_start, c_end))
    }
  }
  return ""
}

var nick = ""
var email = ""

window.onload = function () {
  var conn;
  var msg = document.getElementById("msg");
  var log = document.getElementById("log");

  function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight;
    }
  }

  document.getElementById("form").onsubmit = function () {
    if (!conn) {
      return false;
    }
    if (!msg.value) {
      return false;
    }
    conn.send(msg.value);
    msg.value = "";
    return false;
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/chatsocket");

    // 在这里对连接进行初始化，从cookie中读取personal info，向服务器请求分配资源
    email = getCookie('email')
    nick = getCookie('nick')
    conn.onopen = () => conn.send(nick + "|" + email)

    // 处理是否 remember me
    var r = getCookie('remember')
    if (r == null || r == "") {
      setCookie('email', "", -1)
      setCookie('nick', "", -1)
    }

    conn.onclose = function (evt) {
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };

    // 在这里拆开数据，msg交换协议见OOAD文件夹下的协议文件
    // 注意这里是处理了多条信息，按照后台代码，为了提高负载，当在信道中有
    // 消息开始排队时，后台会一次性的将处于等待状态的消息汇集在一个ws数据包里
    // 以’\n‘分隔
    conn.onmessage = function (evt) {
      var messages = evt.data.split('\n');
      for (var i = 0; i < messages.length; i++) {
        var subMsgs = messages[i].split('|')
        var item = document.createElement("div");
        item.innerText = "***" + subMsgs[2] + "*** [" + subMsgs[0] + "@*.*." + subMsgs[1] + "] "
        for (var j = 3; j < subMsgs.length; j++) {
          item.innerText = item.innerText + subMsgs[j]
        }
        appendLog(item);
      }
    };
  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }
};
