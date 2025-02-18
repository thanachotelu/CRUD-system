import React, { useState, useEffect } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";

const API_URL = "http://localhost:8080/api/v1/user";

const EditUser = () => {
  const { userId } = useParams();
  const navigate = useNavigate();
  const [user, setUser] = useState({
    username: "",
    firstname: "",
    lastname: "",
    email: "",
  });

  useEffect(() => {
    // ดึงข้อมูล user ที่จะทำการแก้ไข
    axios.get(`${API_URL}${userId}`)
      .then(response => {
        setUser(response.data);
      })
      .catch(error => {
        console.error("There was an error fetching the user:", error);
      });
  }, [userId]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser({
      ...user,
      [name]: value
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    axios.put(`${API_URL}${userId}`, user)
      .then(response => {
        // ไปที่หน้าหลักหลังจากการแก้ไขสำเร็จ
        navigate("/");
      })
      .catch(error => {
        console.error("There was an error updating the user:", error);
      });
  };

  return (
    <div>
      <h2>Edit User</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Username</label>
          <input
            type="text"
            className="form-control"
            name="username"
            value={user.username}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>First Name</label>
          <input
            type="text"
            className="form-control"
            name="firstname"
            value={user.firstname}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>Last Name</label>
          <input
            type="text"
            className="form-control"
            name="lastname"
            value={user.lastname}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>Phonenumber</label>
          <input
            type="text"
            className="form-control"
            name="phonenumber"
            value={user.phonenumber}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            className="form-control"
            name="email"
            value={user.email}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>Role</label>
          <input
            type="text"
            className="form-control"
            name="role"
            value={user.role}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label>Status</label>
          <input
            type="text"
            className="form-control"
            name="status"
            value={user.status}
            onChange={handleChange}
          />
        </div>
        <button type="submit" className="btn btn-primary mt-3">
          Update User
        </button>
      </form>
    </div>
  );
};

export default EditUser;
