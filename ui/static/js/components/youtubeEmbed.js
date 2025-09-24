export class YoutubeEmbed extends HTMLElement {
  constructor() {
    super();
    this.startTime = parseInt(this.getAttribute("start") || "0", 10);
    this.initialLoad = false;
    this.playerReady = false;
  }

  connectedCallback() {
    //console.log("youtube embed connected");
    this.embedDiv = this.querySelector("#watchNow");
    this.toggleButtons = document.getElementsByClassName("toggleSketch");
    this.watchNowIframe = this.querySelector("#watchNowIframe");

    if (!(this.embedDiv && this.toggleButtons && this.watchNowIframe)) {
      console.error("Missing elements:", {
        embedDiv: !!this.embedDiv,
        toggleButtons: !!this.toggleButtons && this.toggleButtons.length > 0,
        watchNowIframe: !!this.watchNowIframe,
      });
      throw Error("embed error");
    }

    // Add toggle logic
    for (let button of this.toggleButtons) {
      button.addEventListener("click", () => {
        //console.log("Toggle button clicked");
        this.embedDiv.classList.toggle("hidden");
        this.embedDiv.classList.toggle("flex");

        // Only try to control player if it's ready
        if (this.player && this.playerReady) {
          if (!this.initialLoad) {
            this.player.seekTo(this.startTime, true);
            this.player.playVideo();
            this.initialLoad = true;
          } else {
            if (this.embedDiv.classList.contains("hidden")) {
              this.player.pauseVideo();
            }
          }
        } else if (this.player) {
          //console.log("Player exists but not ready, waiting...");
          this.waitForPlayerReady().then(() => {
            //console.log("Player ready, seeking to:", this.startTime);
            this.player.seekTo(this.startTime, true);
            this.player.playVideo();
          });
        } else {
          console.warn("Player not initialized yet");
        }
      });
    }

    // Wait for API and init
    waitForYouTubeAPI().then(() => {
      //console.log("YouTube API ready, initializing player...");
      this.initPlayer();
    });
  }

  initPlayer() {
    //console.log("Initializing player with iframe:", this.watchNowIframe);

    // Make sure iframe has the required parameters
    const src = this.watchNowIframe.src;
    if (!src.includes("enablejsapi=1")) {
      console.warn("enablejsapi=1 not found in iframe src");
    }

    try {
      // Use the actual iframe element
      this.player = new YT.Player(this.watchNowIframe, {
        events: {
          onReady: (event) => {
            //console.log("YouTube player is ready!", event);
            this.playerReady = true;
          },
          onError: (e) => {
            console.error("YT Player error:", e);
          },
          //onStateChange: (event) => {
          //  console.log("Player state changed:", event.data);
          //},
        },
      });
      //console.log("Player initialization called");
    } catch (error) {
      console.error("Error initializing player:", error);
    }
  }

  waitForPlayerReady() {
    return new Promise((resolve) => {
      if (this.playerReady) {
        resolve();
      } else {
        const checkReady = setInterval(() => {
          if (this.playerReady) {
            clearInterval(checkReady);
            resolve();
          }
        }, 100);

        // Timeout after 10 seconds
        setTimeout(() => {
          clearInterval(checkReady);
          console.warn("Timeout waiting for player to be ready");
          resolve();
        }, 10000);
      }
    });
  }
}

function waitForYouTubeAPI() {
  return new Promise((resolve) => {
    if (window.YT && window.YT.Player) {
      //console.log("YT API already loaded");
      resolve();
    } else {
      //console.log("Waiting for YT API...");
      const check = setInterval(() => {
        if (window.YT && window.YT.Player) {
          //console.log("YT API resolved");
          clearInterval(check);
          resolve();
        }
      }, 100);
    }
  });
}

if (!customElements.get("youtube-embed")) {
  customElements.define("youtube-embed", YoutubeEmbed);
}
