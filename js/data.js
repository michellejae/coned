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
  tooltip: {
    formatter: function (params) {
      let index = params.dataIndex
      return `Name: ${params.name}<br />
      Total: ${params.value}<br />
      Rate: ${graphData[index].rate}
      `
    }
  },
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
  let name = graphData[i].name
  let total = graphData[i].total
  // let minTerm = graphData[i].minTerm
  // let offerType = graphData[i].offerType
  // let rate = graphData[i].rate
  // let cancellation = graphData[i].cancellation
  // let energySource = graphData[i].energySource
  // let percentRenew = graphData[i].percentRenew

  let seriesDataObj = {}
  seriesDataObj.name = name
  seriesDataObj.value = total
  seriesDataObj.itemStyle = {}
  seriesDataObj.itemStyle.color = "blue"

  seriesDataObj.name === "Consolidated Edison Company of New York, Inc." ? seriesDataObj.itemStyle.color = "green" : seriesDataObj.itemStyle.color = "blue"
  
  //option.tooltip.formatter = `${name} <br/> ${total}`
  
  option.series[0].data.push(seriesDataObj)
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


