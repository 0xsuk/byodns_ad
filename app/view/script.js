const init = (path) => {
  let p = document.querySelector("a[href='/" + path + "']");
  p.style.borderBottom = "solid 2px rgb(244, 57, 97)";

  switch (path) {
    case "dashboard":
      break;
    default:
      fetch("/api/" + path + "/read")
        .then((resp) => resp.json())
        .then((resp) => update_list(resp, path));
      console.log("init", path, "page");
  }
};
const set_search_event = (path) => {
  console.log("setting event for", path);
  let elem = document.getElementById("bd-search");
  elem.onkeyup = (e) => {
    let data = elem.value;
    console.log("searching", data);
    let status = document.getElementById("bd-status");
    fetch("/api/" + path + "/search?value=" + data)
      .then((resp) => resp.json())
      .then((resp) => update_list(resp, path))
      .then(() => {
        status.innerText = "";
      })
      .catch((err) => {
        console.log(err);
      });
    status.innerText = "searching";
  };
};
const set_addBlacklist_event = () => {
  console.log("set addBlacklist event");
  let elems = document.getElementsByClassName("addBlacklist");
  elems = Array.from(elems);
  elems.forEach((elem) => {
    elem.onclick = (e) => {
      console.log("adding Blacklist");
      let domain = e.target
        .closest(".domain")
        .querySelector("p:nth-child(1)").innerText;
      let formData = new FormData();
      formData.append("value", domain);
      fetch("/api/blacklist/create", {
        method: "POST",
        body: formData,
      })
        .then((resp) => {
          if (resp.status == 200) {
            elem.className = "addWhitelist";
            elem.innerText = "Whitelist";
          }
        })
        .catch((err) => {
          console.log(err);
        });
      set_addWhitelist_event();
    };
  });
};
const set_addWhitelist_event = () => {
  console.log("set addWhitelist event");
  let elems = document.getElementsByClassName("addWhitelist");
  elems = Array.from(elems);
  elems.forEach((elem) => {
    elem.onclick = (e) => {
      console.log("adding Whitelist");
      let domain = e.target
        .closest(".row")
        .querySelector(".domain p:nth-child(1)").innerText;
      let formData = new FormData();
      formData.append("value", domain);
      fetch("/api/whitelist/create", {
        method: "POST",
        body: formData,
      })
        .then((resp) => {
          if (resp.status == 200) {
            elem.className = "addBlacklist";
            elem.innerText = "Blacklist";
          }
        })
        .catch((err) => {
          console.log(err);
        });
      set_addBlacklist_event();
    };
  });
};
const set_editBlacklist_event = () => {
  console.log("setting editBlacklist event");
  let elems = document.getElementsByClassName("editBlacklist");
  elems = Array.from(elems);
  elems.forEach((elem) => {
    elem.onclick = (e) => {
      console.log("editing blacklist");
      let domain = e.target
        .closest(".row")
        .querySelector(".domain p:nth-child(1)").innerText;
      let newdomain = prompt("type new domain:", domain);
      if (newdomain === null) return;
      let formData = new FormData();
      formData.append("value", newdomain);
      formData.append("old", domain);
      fetch("/api/blacklist/update", {
        method: "PUT",
        body: formData,
      })
        .then((resp) => resp.json())
        .then((resp) => update_list(resp, "blacklist"));
    };
  });
};
const set_delBlacklist_event = () => {
  console.log("setting delBlacklist event");
  let elems = document.getElementsByClassName("delBlacklist");
  elems = Array.from(elems);
  elems.forEach((elem) => {
    elem.onclick = (e) => {
      console.log("deleting blacklist");
      let domain = e.target
        .closest(".row")
        .querySelector(".domain p:nth-child(1)").innerText;
      let formData = new FormData();
      formData.append("value", domain);
      fetch("/api/blacklist/delete", {
        method: "DELETE",
        body: formData,
      }).then((resp) => {
        if (resp.status == 200) {
          elem.className = "";
          elem.innerText = "Removed!";
        }
      });
    };
  });
};

const update_list = (resp, path) => {
  console.log("udpate", path, "page");
  let add_leading_zero = (time) => {
    return ("0" + time).slice(-2);
  };
  let list = document.getElementById("bd-list");
  while (list.firstChild) {
    list.removeChild(list.firstChild);
  }
  switch (path) {
    case "query":
      resp.forEach((query) => {
        let row = document.createElement("div");
        row.className = "row";
        let domain = document.createElement("div");
        domain.className = "domain";

        let time = new Date(query.Timestamp);
        let year = time.getFullYear();
        let month = add_leading_zero(time.getMonth() + 1);
        let date = add_leading_zero(time.getDate());
        let hours = add_leading_zero(time.getHours());
        let minutes = add_leading_zero(time.getMinutes());
        let seconds = add_leading_zero(time.getSeconds());
        let millisec = add_leading_zero(time.getMilliseconds());

        let p1 = document.createElement("p");
        let p2 = document.createElement("p");
        let p3 = document.createElement("p");
        let p4 = document.createElement("p");
        let p5 = document.createElement("p");
        p1.innerText = query.Domain;
        p2.innerText = query.Status;
        p3.innerText = query.ClientIP;
        p4.innerText =
          year +
          "/" +
          month +
          "/" +
          date +
          " " +
          hours +
          ":" +
          minutes +
          ":" +
          seconds +
          "." +
          millisec;
        p5.innerText = query.IsBlocked == "no" ? "Blacklist" : "Whitelist";
        p5.className =
          query.IsBlocked == "no" ? "addBlacklist" : "addWhitelist";

        if (query.Domain == query.OrganizerDomain) {
          row.style.backgroundColor = "#090528";
        }

        domain.appendChild(p1);
        domain.appendChild(p2);
        domain.appendChild(p3);
        domain.appendChild(p4);
        domain.appendChild(p5);
        row.appendChild(domain);
        list.appendChild(row);
      });
      //FIX of "object ..."
      set_addBlacklist_event();
      set_addWhitelist_event();
      break;
    case "blacklist" || "whitelist":
      resp.forEach((elem) => {
        let row = document.createElement("div");
        row.className = "row";
        let domain = document.createElement("div");
        domain.className = "domain";
        let edit = document.createElement("div");
        edit.className = "editBlacklist";
        let del = document.createElement("div");
        del.className = "delBlacklist";
        let p1 = document.createElement("p");
        p1.innerText = elem;
        let p2 = document.createElement("p");
        p2.innerText = "Edit";
        let p3 = document.createElement("p");
        p3.innerText = "Delete";

        domain.appendChild(p1);
        edit.appendChild(p2);
        del.appendChild(p3);
        row.appendChild(domain);
        row.appendChild(edit);
        row.appendChild(del);
        list.appendChild(row);
      });

      set_editBlacklist_event();
      set_delBlacklist_event();
      break;
  }
};

//should be called before window.onload
let path = document.currentScript.getAttribute("path");
window.onload = () => {
  init(path);
  if (path == "blacklist" || path == "whitelist" || path == "query") {
    set_search_event(path);
  }
  if (path == "query") {
    set_addBlacklist_event();
  }
};
