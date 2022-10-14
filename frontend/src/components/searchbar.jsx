import React from "react";
import { useEffect, useCallback } from "react";
import { Search } from "semantic-ui-react";
import axios from "axios";
import { useState } from "react";

export const Searchbar = ({
  onResultSelected
}) => {
  const [cities, setCities] = useState([]);
  const [reqInprogress, setReqInprogress] = useState(false);
  const [filteredCities, setFilteredCities] = useState([]);

  useEffect(() => {
    if (!reqInprogress) {
      (async () => {
        setReqInprogress(true);
        const citiesResponse = await axios.get(
          "http://localhost:5001/cities/search"
        );

        const formattedCities = Object.values(citiesResponse.data || {}).map(
          ({ id, name }) => ({
            title: id,
            description: name,
          })
        );
        setCities(formattedCities);
      })();
    }
  }, [reqInprogress]);

  const handleSearchChange = useCallback(
    (_, data) => {
      const searchString = data.value.toLowerCase();
      const filterdCities = cities.filter(
        (i) =>
          i.title.toLowerCase().indexOf(searchString) > -1 ||
          i.description.toLowerCase().indexOf(searchString) > -1
      );
      setFilteredCities(filterdCities);
    },
    [cities]
  );
  return (
    <div className="search-bar">
      <Search
        loading={false}
        placeholder="Start journey from.."
        onResultSelect={(e, data) => {
          onResultSelected(data.result.title);
        }}
        onSearchChange={handleSearchChange}
        results={filteredCities || cities}
        size="big"
        minCharacters={2}
        fluid={true}
      />
    </div>
  );
};
