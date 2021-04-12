import React, { Component } from 'react';
import './Footer.styles.css';
class Footer extends Component {
  state = {};
  render() {
    return (
      <div className="footer">
        <div>
          <h2>Bitly - CMPE 281</h2>
          <h3>Under Prof. Paul Nguyen</h3>
        </div>
        <ul>
          <h4>Developer</h4>
          <li>Vishal Gadapa</li>
        </ul>
        <p>Icons from Bilty website</p>
      </div>
    );
  }
}

export default Footer;
