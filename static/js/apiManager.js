// get api
function getApi(url, callback) {
  fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then(function (response) {
      console.log(response);
      return response.json();
    })
    .then(function (data) {
      console.log(data);
      callback(data); // 응답에 성공한 후 콜백 함수 실행
    })
    .catch(function (error) {
      console.log("ajax error:", error);
    });
}

// post
function postApi(url, data, callback) {
  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (!response.ok) {
        console.log(response);
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((responseData) => {
      console.log(responseData);
      if (callback) {
        callback(responseData); // 응답에 성공한 후 콜백 함수 실행
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}
