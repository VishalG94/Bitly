module.exports = {
  BACKEND_SERVER: {
    URL: "http://34.221.223.237:8000/cp/",
    // URL: "http://202BackendLB-1797169985.us-east-1.elb.amazonaws.com" // Should have http://
  },
  TREND_SERVER: {
    URL: "http://34.221.223.237:8000/ts/",
    // URL: "http://202BackendLB-1797169985.us-east-1.elb.amazonaws.com" // Should have http://
  },
  USER_INFORMATION: {
    USER_ID: localStorage.getItem("userId"),
    USERNAME: localStorage.getItem("userName"),
    USER_TYPE: localStorage.getItem("userType"),
  },
};
