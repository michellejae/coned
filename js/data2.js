fetch("/", {
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    },
   method: "POST",
}).then(response => response.json())
.then(data => console.log(data, "HERE"));