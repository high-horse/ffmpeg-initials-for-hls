const URL = "http://localhost:9090";

$(document).ready(function () {
  fetchAllSamples();

  $("#uploadForm").on("submit", function (e) {
    e.preventDefault();
    submitSample();
  });

  $("#media-tab").click(function (e) {
    e.preventDefault();
    fetchAllSamples();
  });
});

function fetchAllSamples() {
  $.ajax({
    url: URL + "/hls/all",
    method: "GET",
    success: function (response) {
      console.log("Samples fetched:", response);

      let renderDiv = $("#mediaList");
      renderDiv.empty();
      if (response.length === 0) {
        renderDiv.append("<p>No media found.</p>");
        return;
      }
      // Render a list of media items
      response.forEach((item) => {
        const mediaButton = $(`
               <button class="btn btn-outline-primary mb-2 d-block w-100 text-start" data-path="${item.path}">
                 ▶️ ${item.name}
               </button>
             `);

        // On click, play the video
        mediaButton.on("click", function () {
          const path = $(this).data("path");
          fetchHlsData(path); // assumes fetchHlsData accepts the path directly
        });

        renderDiv.append(mediaButton);
      });
    },
    error: function (xhr, status, error) {
      console.error("GET failed:", error);
    },
  });
}

function submitSample() {
  const file = $("#mediaFile")[0].files[0];
  if (!file) {
    alert("Please choose a file to upload.");
    return;
  }

  const formData = new FormData();
  formData.append("media_file", file);

  $.ajax({
    url: "http://localhost:9090/upload", // your Go server endpoint
    method: "POST",
    data: formData,
    processData: false,
    contentType: false,
    success: function (response) {
      alert("Upload successful!");
      // Activate the media tab
      const mediaTab = new bootstrap.Tab(document.querySelector("#media-tab"));
      mediaTab.show();

      fetchAllSamples();
    },
    error: function (xhr, status, error) {
      alert("Upload failed.");
      console.error(error);
    },
  });
}

function fetchSample(id) {
  $.ajax({
    url: URL + "/hls/sample_id",
    method: "GET",
    contentType: "application/json",
    // data: JSON.stringify({ action: "play", videoId: 123 }),
    success: function (response) {
      console.log("Action sent:", response);
    },
    error: function (xhr, status, error) {
      console.error("POST failed:", error);
    },
  });
}

function fetchHlsData(path) {
  const video = document.getElementById("video");

  // Change this URL to match your Go server
  const hlsUrl = "http://localhost:9090/hls/" + path;

  if (Hls.isSupported()) {
    const hls = new Hls();
    hls.loadSource(hlsUrl);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play();
    });
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    // Safari can play HLS natively
    video.src = hlsUrl;
    video.addEventListener("loadedmetadata", function () {
      video.play();
    });
  } else {
    alert("Your browser does not support HLS playback.");
  }
}

function fetchHlsData_(sampleName) {
  const video = document.getElementById("video");

  // Change this URL to match your Go server
  const hlsUrl = "http://localhost:9090/hls/sample-hls-op.m3u8";

  if (Hls.isSupported()) {
    const hls = new Hls();
    hls.loadSource(hlsUrl);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play();
    });
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    // Safari can play HLS natively
    video.src = hlsUrl;
    video.addEventListener("loadedmetadata", function () {
      video.play();
    });
  } else {
    alert("Your browser does not support HLS playback.");
  }
}
