(async (packet) => {
  
  // entry point
  main = () => {
    switch(packet.stateManager.currState) {
      case "init":
        // set title and heading
        break;
      case "ready":
        // begin pulling route info
        break;
      case "closing":
        // make sure everything is exported accordingly
        break;
      default: 
        // you be fail if find self nar
        break;
    }
  }
  // start
  main()
})
({
  stateManager : {
    currState: "init",
    availStates: [
      "init", "ready", "closing",
    ]
  },
  api : {
    nextrip : {
      format: "format=JSON",
      proxy: "https://cors-anywhere.herokuapp.com",
      nextrip: "https://svc.metrotransit.org",
      api: "nextrip",
      get GetAPIURL() {
        return `${this.proxy}/${this.nextrip}/${this.api}/`
      },
    },
    get GetTodaysRoutes() {
      return new Promise((resolve, reject) => {
        fetch(this.nextrip.GetAPIURL + "Routes" + this.api.format)
          .then((response) => response.json())
          .then((data) => resolve(data))
      })
    },
  },
  presentation: {
    title: "Metro Transit NexTrip API",
    tagline: "repeteable code,"
  }
}) 
