const arr = []

fetch("/", {
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    },
   method: "POST",
}).then(response => response.json())
.then(data => {
  arr.push(data)
});

var myChart = echarts.init(document.getElementById('main'));

console.log(arr, "herrr")
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
    data: ['Shirts', 'Cardigans', 'Chiffons', 'Pants', 'Heels', 'Socks']
  },
  yAxis: {},
  series: [
    {
      name: 'sales',
      type: 'bar',
      data: [5, 20, 36, 10, 10, 20]
    }
  ]
};

// Display the chart using the configuration items and data just specified.
myChart.setOption(option);