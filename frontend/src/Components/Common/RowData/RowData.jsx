import React from "react";
import "./RowData-styles.css";

const RowData = ({ URL, ShortLink, Count }) => {
  return (
    <tr>
      <td href={URL}>{URL}</td>
      <td><a href={ShortLink}>{ShortLink}</a></td>
      <td>{Count}</td>
    </tr>
  );
};

export default RowData;
