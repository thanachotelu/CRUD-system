import React from "react";
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Adduser from './pages/Adduser'
import Edituser from './pages/Edituser'

function App() {
  return (
    <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/adduser" element={<Adduser/>}/>
          <Route path="/edituser/:userId" element={<Edituser/>}/>
        </Routes>
      </Router>
  );
}

export default App;
