import React, { Component } from "react";
import "./LandingPage.styles.css";

import Jumbotron from "../../Components/Common/Jumbotron/Jumbotron";
import Banner from "../Common/Banner/Banner";
import howData from "../Common/BannerDataItems/How";
import whyData from "../Common/BannerDataItems/Why";

import Constants from "../../Utils/Constants";
import axios from "axios";
import RowData from "../Common/RowData/RowData";

class LandingPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      trendingLinks:[]
    };
  }
  componentDidMount() {
    axios.get(`${Constants.TREND_SERVER.URL}/shortlinktrend`)
    .then((response) => {
      console.log(response.data);
      this.setState({ trendingLinks: response.data })
    }).catch((error) => {
        console.log(error)
    });
  }
  render() {
    const { trendingLinks } = this.state;
    return (
      <div className="homeLayout">
        <Jumbotron />
        <h2>Trending links! </h2>
        <div className="vehicleCatalog">
        </div>
        <div >
        <table class="table table-hover">
  <thead>
    <tr>
      <th scope="col">URL</th>
      <th scope="col">Short Link</th>
      <th scope="col">Count</th>
    </tr>
  </thead>
  <tbody>
  {trendingLinks.map((details) => (
                <RowData {...details} />
  ))} 
  </tbody>
</table>
      </div>
      </div>
    );
  }
}

export default LandingPage;
