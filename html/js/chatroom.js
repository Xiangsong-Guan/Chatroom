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
    // TODO 在这里处理数据，msg交换协议见OOAD文件夹下的协议文件
    // 像是数据处理这种任务，能放在客户端的都放在客户端，省服务器的计算资源
    conn.send(msg.value);
    msg.value = "";
    return false;
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/chatsocket");
    conn.onclose = function (evt) {
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };
    conn.onmessage = function (evt) {
      var messages = evt.data.split('\n');
      for (var i = 0; i < messages.length; i++) {
        // TODO 在这里拆开数据，msg交换协议见OOAD文件夹下的协议文件
        // 注意这里是处理了多条信息，按照后台代码，为了提高负载，当在信道中有
        // 消息开始排队时，后台会一次性的将处于等待状态的消息汇集在一个ws数据包里
        // 以’\n‘分隔
        var item = document.createElement("div");
        item.innerText = messages[i];
        appendLog(item);
      }
    };
  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }

  // 在这里对连接进行初始化，从cookie中读取personal info，向服务器请求分配资源
};
