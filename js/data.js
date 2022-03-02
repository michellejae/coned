const myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");


const requestOptions = {
  method: "POST",
  headers: myHeaders,
};

async function foobar(){
  let graphData;

  try {
    const response = await fetch("/", requestOptions)
    graphData = await response.json()
  } catch (error) {
    console.log(error)
  }


var myChart = echarts.init(document.getElementById('main'));


// Specify the configuration items and data for the chart
var option = {
  title: {
    text: 'ECharts Getting Started Example'
  },
  tooltip: {},
  legend: {
  },
  xAxis: {
    type: "category",
    data: [],
    show: false
  },
  yAxis: {},
  series: [
    {
      type: 'bar',
      data: [],
    }
  ]
};





for (let i=0 ; i<graphData.length; i++){
  let obj = {}
  obj.name = graphData[i].name
  obj.value = graphData[i].total
  obj.itemStyle = {}
  obj.itemStyle.color = "blue"

  obj.name === "Consolidated Edison Company of New York, Inc." ? obj.itemStyle.color = "green" : obj.itemStyle.color = "blue"
  
  //option.xAxis.data.push(graphData[i].name)
  
  option.series[0].data.push(obj)
}




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


