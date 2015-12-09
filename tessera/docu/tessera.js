window.addEventListener('load', function(){

  document.getElementById("login").style.visibility = "collapse";
  var el = document.getElementById("regist");
  var log = document.getElementById("log");
  var data = new FormData();
  var xhr = new XMLHttpRequest();

  el.addEventListener('click', function(){

    var pass1 = document.getElementById("pass1").value;
    var pass2 = document.getElementById("pass2").value;
    var name = document.getElementById("user").value;
    data = document.getElementById("reg");
    if (name != "") {
      if(pass1 == pass2)
      {
        document.getElementById("text").innerHTML = "Willkommen "+name + "!";
		      xhr.open( 'GET', 'http://localhost:4242/regist', true );
		      xhr.send(data);
          document.getElementById("test").innerHTML = data + "";
      }
      else
      {
        document.getElementById("text").innerHTML = "Angegebene Passwörter stimmen nicht überein";
      }
    }
    else
    {
      document.getElementById("text").innerHTML = "Bitte noch einmal ausfüllen!";
    }

  });

  var elem = document.getElementById("login");
  var elem1 = document.getElementById("text");
  var elem2 = document.getElementById("reg");

  log.addEventListener("click", function(){

    if (elem.style.visibility =="collapse")
    {
      elem.style.visibility = "visible";
      elem1.innerHTML = "";
      elem2.style.visibility = "collapse";
    }
  });

  log2.addEventListener("click", function(){
    if (elem2.style.visibility =="collapse")
    {
      elem2.style.visibility = "visible";
      elem1.innerHTML = "";
      elem.style.visibility = "collapse";
    }
  });
});
