import React from "react";

import airplane from "../assets/airplane.png";

export const Navbar = () => {
  return (
    <nav>
      <div className="logo">
      <img src={airplane} height="40" width="40" alt="Airplane Logo" />
      <span className="logo-text">
        Globe Trip
      </span>
      </div>
    </nav>
  );
};
