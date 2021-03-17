const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    var first = true; 

    for (let result of results) {
      if (first){ 
          var n = result.search(data.query);
          var part1 = result.substring(0,n);
          var l = length(data.query);
          var part2 = result.substring(n+l);
          result = `${part1}<mark>${data.query}</mark>${part2}`;
          first = false;
        }
      rows.push(`<tr>${result}<tr/>`);
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
