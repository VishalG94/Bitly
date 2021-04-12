import React from "react";
import "./LinkTag-style.css";
const LinkTag = ({ handleChange, label, ...otherProps }) => (
  <div className="group">
    <input className="form-input" onClick={handleChange} {...otherProps} />
    {label ? (
      <label
        className={`${
          otherProps.value.length ? "shrink" : ""
        } form-input-label`}
      >
        {label}
      </label>
    ) : null}
  </div>
);

export default LinkTag;
