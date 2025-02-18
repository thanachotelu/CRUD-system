import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const API_URL = "http://localhost:8080/api/v1/user";

const AddUser = () => {
  const navigate = useNavigate();
  const [newUser, setNewUser] = useState({
    username: "",
    password: "",
    firstname: "",
    lastname: "",
    phonenumber: "",
    email: "",
    role: "",
    status: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewUser((prevUser) => ({
      ...prevUser,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post(API_URL, newUser);
      navigate("/"); // กลับไปหน้าหลักหลังจากเพิ่มข้อมูลเสร็จ
    } catch (error) {
      console.error("Error adding user:", error);
    }
  };

  return (
    <div className="container mt-5">
      <h1 className="text-center mb-4">Add New User</h1>
      <form onSubmit={handleSubmit}>
        <div className="mb-3">
          <label>Username</label>
          <input
            type="text"
            className="form-control"
            name="username"
            value={newUser.username}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Password</label>
          <input
            type="password"
            className="form-control"
            name="password"
            value={newUser.password}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>First Name</label>
          <input
            type="text"
            className="form-control"
            name="firstname"
            value={newUser.firstname}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Last Name</label>
          <input
            type="text"
            className="form-control"
            name="lastname"
            value={newUser.lastname}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Phone Number</label>
          <input
            type="text"
            className="form-control"
            name="phonenumber"
            value={newUser.phonenumber}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Email</label>
          <input
            type="email"
            className="form-control"
            name="email"
            value={newUser.email}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Role</label>
          <input
            type="text"
            className="form-control"
            name="role"
            value={newUser.role}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label>Status</label>
          <input
            type="text"
            className="form-control"
            name="status"
            value={newUser.status}
            onChange={handleChange}
            required
          />
        </div>
        <button type="submit" className="btn btn-primary w-100">
          Add User
        </button>
      </form>
    </div>
  );
};

export default AddUser;
