fetch("/", {
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    },
  //  method: "POST",
    body: JSON.stringify(data)
}).then(response => response.json())
.then(data => console.log(data));