import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './Home.css';

const Home = () => {
  const [key, setKey] = useState('');
  const [ttl, setTTL] = useState(5);
  const [cache, setCache] = useState({});
  const [error, setError] = useState(null);
  const [showRightContainer, setShowRightContainer] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      fetchCacheKeys();
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const fetchCacheKeys = async () => {
    try {
      const response = await axios.get('http://localhost:8080/cache');
      const cacheData = response.data.cache;

      setCache(cacheData);
      setShowRightContainer(Object.keys(cacheData).length > 0);
    } catch (error) {
      console.error('Error fetching cache keys:', error);
    }
  };

  const resetInputs = () => {
    setKey('');
    setTTL(5);
  };

  const getCache = async (key) => {
    try {
      const response = await axios.get(`http://localhost:8080/cache/${key}`);
      setError(null);
      resetInputs();
      setKey('');
    } catch (error) {
      setError('Key not found');
      resetInputs();
    }
  };

  const setCacheItem = async () => {
    const ttlValue = parseInt(ttl);

    if (key.trim() === '' || isNaN(ttlValue) || ttlValue <= 0) {
      setError('Please enter a valid Key and TTL (seconds)');
      return;
    }

    try {
      await axios.post('http://localhost:8080/cache', { key, ttl: ttlValue });
      fetchCacheKeys();
      setError(null);
      resetInputs();
    } catch (error) {
      setError('Error setting key');
    }
  };

  const deleteCacheItem = async (key) => {
    try {
      await axios.delete(`http://localhost:8080/cache/${key}`);
      fetchCacheKeys();
      setError(null);
    } catch (error) {
      setError('Error deleting key');
    } finally {
      resetInputs();
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    if (name === 'key') setKey(value);
    else if (name === 'ttl') setTTL(value);
  };

  return (
    <div className="home-container">
      <div className="left-container">
        <h1>CacheSwift Operations</h1>
        {error && <p className="error">{error}</p>}
        <div className="input-section">
          <h2>Set Cache Item</h2>
          <input
            type="text"
            name="key"
            placeholder="Key"
            value={key}
            onChange={handleInputChange}
          />
          <input
            type="number"
            name="ttl"
            placeholder="TTL (seconds)"
            value={ttl}
            onChange={handleInputChange}
          />
          <button className="btn" onClick={setCacheItem}>Set</button>
        </div>
        <div className="input-section">
          <h2>Get Cache Item</h2>
          <input
            type="text"
            name="key"
            placeholder="Key"
            value={key}
            onChange={handleInputChange}
          />
          <button className="btn" onClick={() => getCache(key)}>Get</button>
        </div>
        <div className="input-section">
          <h2>Delete Cache Item</h2>
          <input
            type="text"
            name="key"
            placeholder="Key"
            value={key}
            onChange={handleInputChange}
          />
          <button className="btn" onClick={() => deleteCacheItem(key)}>Delete</button>
        </div>
      </div>
      {showRightContainer && (
        <div className="right-container">
          <h2>Current Cache Keys</h2>
          <ul>
            {Object.keys(cache).map((key) => (
              <li key={key}>
                {key}: (Expires at {new Date(cache[key].expiration).toLocaleString()})
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default Home;
