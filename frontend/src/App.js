import { useState } from "react";
import "./App.css";
import { Navbar } from "./components/navbar";
import { TripDetails } from "./components/tripDetails";
import { Searchbar } from "./components/searchbar";

function App() {
  const [selectedCity, setSelectedCity] = useState(undefined);
  const onResultSelected = (city)=>{
    setSelectedCity(city);
  }

  return (
    <div>
      <Navbar />
      <Searchbar onResultSelected={onResultSelected} />
      <TripDetails selectedCity ={selectedCity}/>
    </div>
  );
}

export default App;
