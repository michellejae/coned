const myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");


const requestOptions = {
  method: "POST",
  headers: myHeaders,
};

async function foobar(){
  let data;

  try {
    const response = await fetch("/", requestOptions)
    data = await response.json()
  } catch (error) {
    console.log(error)
  }


var myChart = echarts.init(document.getElementById('main'));

let name = [];
let rate = [];


for (let i=0 ; i<data.length; i++){
  name.push(data[i].name)
  rate.push(data[i].rate)
}



// Specify the configuration items and data for the chart
var option = {
  title: {
    text: 'ECharts Getting Started Example'
  },
  tooltip: {},
  legend: {
    data: ['sales']
  },
  xAxis: {
    data: name
  },
  yAxis: {},
  series: [
    {
      name: 'sales',
      type: 'bar',
      data: rate
    }
  ]
};

// Display the chart using the configuration items and data just specified.
myChart.setOption(option);

}

foobar()

// async function getData(){
//   const response = await fetch("/")
//   const sources = await response.json()
//   return sources
// }


// fetch("/", {
//     headers: {
//         'Accept': 'application/json',
//         'Content-Type': 'application/json'
//     },
//    method: "POST",
// }).then(response => response.json())
// .then(data => {
  
// })


