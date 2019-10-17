const endpoint = "https://cors-anywhere.herokuapp.com/https://svc.metrotransit.org/nextrip"

let routeNumber = "0"
let directional = "0"
let stop = "XXXX"
let stopName = ""
let departuresList = []

let setStatus = (status) => document.querySelector("#status").innerText = status
let setMarquee = (departuresString) => document.querySelector('#marquee-text').innerText = departuresString
let unhide = (className) => document.querySelector(`.${className}`).classList.remove('hidden')
let resetSelects = () => {
  let opt = selRoute.options[0]
  selRoute.innerText = ""
  selRoute.appendChild(opt)
  selDirection.innerText = ""
  selDirection.appendChild(opt)
  selStop.innerText = ""
  selStop.appendChild(opt)
}

let getRoutes = () => {
  
  let routesEndpoint = `${endpoint}/Routes?format=JSON`
  setStatus('fetching routes...')
  fetch(routesEndpoint)
  .then((response) => response.json())
  .then((routes) => {
    setStatus('routes received..')
    routes.forEach((route) => {
      let opt = document.createElement('option')
      let txtNode = document.createTextNode(route.Description)
      opt.value = route.Route
      opt.appendChild(txtNode)
      selRoute.appendChild(opt)
    })
    unhide('routes')
    selRoute.focus()
  })
  .catch((err) => { throw Error(err)} )
}

let routePicked = (event) => {
  let opt = selRoute.options[0]
  
  console.log(event.target.options.length)

  if (selStop.options.lengh > 1) {
    selStop.innerText = ""
    selStop.appendChild(opt)
  }
  let dex = event.target.selectedIndex
  routeNumber = event.target.options[dex].value
  let ep = `${endpoint}/Directions/${routeNumber}?format=JSON`
  setStatus(`${routeNumber} - fetching directions...`)
  fetch(ep)
  .then((response) => response.json())
  .then((directions) => {
    setStatus(`${routeNumber} - directions received...`)
    directions.forEach((direction) => {
      let opt = document.createElement('option')
      let txt = document.createTextNode(direction.Text)
      opt.value = direction.Value
      opt.appendChild(txt)
      selDirection.appendChild(opt)
    })
    unhide('directions')
    selDirection.focus()
  })

}

let directionPicked = (event) => {
  let opt = selRoute.options[0]

  console.log(event.target.options.length)

  let dex = event.target.selectedIndex
  directional = event.target.options[dex].value
  let ep = `${endpoint}/Stops/${routeNumber}/${directional}?format=JSON`
  setStatus(`${routeNumber} - ${directional} - fetching stops...`)
  fetch(ep)
  .then((response) => response.json())
  .then((stops) => {
    setStatus(`${routeNumber} - ${directional} - stops received...`)
    stops.forEach((stop) => {
      let opt = document.createElement('option')
      let txt = document.createTextNode(stop.Text)
      opt.value = stop.Value
      opt.appendChild(txt)
      selStop.appendChild(opt)
    })
    unhide('stops')
    selStop.focus()
  })
}

let stopPicked = (event) => {
  let dex = event.target.selectedIndex;
  stop = event.target.options[dex].value
  let ep = `${endpoint}/${routeNumber}/${directional}/${stop}?format=JSON`
  setStatus(`${routeNumber} - ${directional} - ${stop} - fetching departures...`)
  fetch(ep)
  .then((response) => response.json())
  .then((departures) => {
    setStatus("departures received...")
    departuresList = departures
    let msg = `${routeNumber} - ${directional} - ${stop} - Bus arriving at ${departures[0].DepartureText}`
    
    setStatus(msg)
    unhide('departure-list')
    let departList = document.querySelector("#next-departures-list")
    let tStop = event.target.options[dex].innerText
    let tDescription = departures[0].Description
    let tGate = departures[0].Gate
    let tHeading = departures[0].RouteDirection
    document.querySelector("#departure-list-details").innerText = `Departing ${tStop} ${tDescription} at Gate ${tGate} heading ${tHeading}`

    departures.forEach((depart) => {
      console.log(JSON.stringify(depart,null,2))
      let item = document.createElement('li')
      let text = document.createTextNode(`${depart.DepartureText}`)
      item.appendChild(text)
      departList.appendChild(item)
    })
    
  })

}

let loaded = () => {
  document.querySelector("#title").innerText = "Metro Transit NexTrip API"
  selRoute.addEventListener("change", routePicked)
  selDirection.addEventListener("change", directionPicked)
  selStop.addEventListener("change", stopPicked)
  getRoutes()
}

let selRoute = document.querySelector("#select-route")
let selDirection = document.querySelector("#select-direction")
let selStop = document.querySelector("#select-stop")

loaded()