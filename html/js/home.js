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

window.addEventListener("load", () => {
  email = getCookie('email')
  nick = getCookie('nick')
  last = getCookie('success')
  if (last != null && last != "") {
    tip = document.getElementById('tip')
    tip.innerHTML = '<h1 class="h3 mb-3 font-weight-normal">欢迎回到聊天室</h1><h3 class="h3 mb-3 font-weight-normal" >同志 ' + nick + '</h3>'

    document.getElementById('inputEmail').value = email
    document.getElementById('inputNick').value = nick
    document.getElementById('remember').checked = true
    setTimeout(() => {
      document.getElementById("form").submit()
    }, 1500);
  }
})

document.getElementById("form").addEventListener("submit", () => {
  setCookie('email', document.getElementById('inputEmail').value, 30)
  setCookie('nick', document.getElementById('inputNick').value, 30)
  if (document.getElementById('remember').checked) {
    setCookie('remember', document.getElementById('remember').value, 30)
  } else {
    setCookie('remember', "", -1)
  }
  setCookie('success', "", -1)
})