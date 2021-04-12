import React, { Component } from "react";
import img from "../../../Assets/Images/JumbotronImage.jpg";
import "./Jumbotron.styles.css";
import "../LinkTag/LinkTag"
import Constants from "../../../Utils/Constants";
import axios from "axios";

class Jumbotron extends Component {
  state = {};

  constructor(props) {
    super(props);
    this.state = {
      link: "",
      successMsg: false,
      existingLink:false,
      shortLink: "",
      // submitDiabled:true
    };
  }
  
  handleChange = (e) => {
    const { value, name } = e.target;
    this.setState({ [name]: value });
    this.setState({successMsg:false})
    this.setState({existingLink:false})
  };

  handleSubmit = async (e) => {

    e.preventDefault();
    const userdetails = {
      URL: this.state.link
    };
    console.log(userdetails);
    console.log(Constants.BACKEND_SERVER.URL)
    // localhost:3000/shortlink
    axios
      .post(`${Constants.BACKEND_SERVER.URL}/shortlink`, userdetails)
      .then((response) => {
        console.log(response)
        if (response.status === 200) {
          console.log(JSON.stringify(response.data))
          this.setState({successMsg:true})
          this.setState({shortLink:response.data.ShortLink})
        } else {
          console.log("Reposne is incorrect!")
        }
      }).catch((err) => {
        if(err.response!== undefined && err.response.status === 409){
          console.log("Links: "+ err.response.data)
          this.setState({existingLink:true})
          this.setState({successMsg:false})
          this.setState({shortLink:err.response.data})
        }else{
          console.log("Error Encountered!")
        }
      });
  };

  render() {
    let linkCreationStatus = null
    if(this.state.successMsg){
      linkCreationStatus = 
    <div>
      <p class="successText">ShortLink Created Successfully</p>
      <a href={this.state.shortLink}>{this.state.shortLink}</a>
    </div>
    } 

    let linkCreationFailed = null
    if(this.state.existingLink){
      linkCreationFailed = 
    <div>
      <p class="existingLinkText">Link Already Existing!</p>
      <a href={this.state.shortLink}>{this.state.shortLink}</a>
    </div>
    } 

    return (
      <div class="jumbocontainer">
        <div className="jumbotronImage">
          <img src={img} alt="JumbotronImage"></img>
        </div>
        <br/>
        <div className="jumbotronText">
          <p class="title">Short links, big results</p>
          <p class="subtext">
          {/* Short links, big results */}
          A URL shortener built with powerful tools to help you grow and protect your brand.
          </p>
          <br>
          </br>
          <div class="input-group mb-3">
            <input type="text" class="form-control" placeholder="Complete Link" aria-label="Complete Link" aria-describedby="basic-addon2" name="link" onChange={this.handleChange} />
            <div class="input-group-append">
              <button class="btn btn-outline-primary" type="button" onClick={this.handleSubmit} disabled={!this.state.link}>Shorten</button>
            </div>
          </div>
          {linkCreationStatus}
          {linkCreationFailed}
        </div>
      </div>
    );
  }
}

export default Jumbotron;
