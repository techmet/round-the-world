import axios from "axios";
import { useState } from "react";
import { useEffect } from "react";

export const TripDetails = ({ selectedCity }) => {
  const [tripDetails, setTripDetails] = useState(undefined);
  useEffect(() => {
    if (selectedCity) {
      (async () => {
        const tripDetails = await axios.get(
          `http://localhost:5001/cities/trip/${selectedCity}`
        );
        setTripDetails(tripDetails.data);
      })();
    }
  }, [selectedCity]);

  if (tripDetails) {
    return (
      <div class="trip-details">
        <div class="trip-details-header">Trip Details</div>
        <div class="trip-details-content">
          <b>Trip:</b> {tripDetails.path.join(" ➙ ")}
        </div>
        <div class="trip-details-content">
          <b>Distance Travelled</b>: {Math.round(tripDetails.totalDistance)} Kms
        </div>
      </div>
    );
  }

  return null;
};
