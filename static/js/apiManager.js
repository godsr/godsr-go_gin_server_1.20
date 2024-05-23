/*
  Axios 라이브러리 사용 중
  사용 예시
      axiosRequest(method, url, requestData)
      .then(function (data) {
        console.log("Success:", data);
        성공시 사용할 함수(data);
      })
      .catch(function (errMsgData) {
        alert(errMsgData.result);
      });

*/

/*
  Axios interceptor를 사용하여 request에 access token 삽입
*/
axios.interceptors.request.use(
  function (config) {
    const accessToken = localStorage.getItem("accessToken");
    if (accessToken) {
      config.headers["Authorization"] = `Bearer ${accessToken}`;
    }
    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

/*
  access token이 만료시 refresh token으로 access token 재발급을 요청
*/
async function refreshAccessToken() {
  try {
    const refreshToken = localStorage.getItem("refreshToken");
    const response = await axios.post(
      "http://localhost:8080/api/test/refresh",
      {
        refreshToken: refreshToken,
      }
    );
    const newAccessToken = response.data.accessToken;
    const newRefreshToken = response.data.refreshToken;
    localStorage.setItem("accessToken", newAccessToken);
    localStorage.setItem("refreshToken", newRefreshToken);
    return {
      accessToken: newAccessToken,
      refreshToken: newRefreshToken,
    };
  } catch (error) {
    throw new Error("Failed to refresh access token");
  }
}

/*
  Axios 를 사용하여 Http 요청
*/

function axiosRequest(method, url, data = null) {
  // Axios 요청 설정
  const config = {
    method: method,
    url: url,
    data: data,
  };

  // Axios 요청을 Promise로 반환
  return axios(config)
    .then(function (response) {
      // 성공 시 response.data 반환
      return response.data;
    })
    .catch(async function (error) {
      // 엑세스 토큰이 만료된 경우
      if (error.response && error.response.status === 401) {
        try {
          const newToken = await refreshAccessToken();
          // 새로운 엑세스 토큰으로 원래 요청을 다시 시도
          config.headers = {
            ...config.headers,
            Authorization: `Bearer ${newToken.accessToken}`,
          };
          if ("refreshToken" in config.data) {
            config.data.refreshToken = newToken.refreshToken;
          }
          const retryResponse = await axios(config);
          return retryResponse.data;
        } catch (refreshError) {
          return Promise.reject(refreshError.message);
        }
      } else {
        // 실패 시 에러 메시지를 reject
        if (error.response && error.response.data) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject("An error occurred: " + error.message);
        }
      }
    });
}
