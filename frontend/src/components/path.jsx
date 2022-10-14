import axios from "axios";
import { useState } from "react";
import { useEffect } from "react";

export const Path = ({ selectedCity }) => {
  const [tripDetails, setTripDetails] = useState(undefined);
  useEffect(() => {
    if (selectedCity) {
      (async () => {
        const tripDetails = await axios.get(
          `http://localhost:5001/cities/trip/${selectedCity}`
        );
        tripDetails.data.path.unshift(tripDetails.data.path.at(-1))
        setTripDetails(tripDetails.data);
      })();
    }
  }, [selectedCity]);

  if (tripDetails) {
    return (
      <div class="trip-details">
        <div class="trip-details-header">Trip Details</div>
        <div class="trip-details-content">
          <b>Total Distance</b>: {Math.round(tripDetails.totalDistance)} 
        </div>
        <div class="trip-details-content">
       <b>Trip:</b> {tripDetails.path.join("/ ")}
        </div>
      </div>
    );
  }

  return null;
};
