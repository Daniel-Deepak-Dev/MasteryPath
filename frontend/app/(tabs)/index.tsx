import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TextInput,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
  Alert,
  SafeAreaView,
} from 'react-native';
import axios from 'axios';

// IMPORTANT: Change this IP address to match your setup
// For Android Emulator: use 10.0.2.2
// For iOS Simulator: use localhost
// For physical device: use your computer's local IP (e.g., 192.168.1.100)
const BASE_API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api';
const GOALS_API_URL = `${BASE_API_URL}/goals`;

// TypeScript interface for Goal
interface Goal {
  id: string;
  title: string;
  description: string;
  completed: boolean;
  created_at: string;
}

export default function MasteryPathScreen() {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null);

  // Fetch goals when the screen loads
  useEffect(() => {
    fetchGoals();
  }, []);

  // GET: Fetch all goals from the backend
  const fetchGoals = async () => {
    setLoading(true);
    try {
      const response = await axios.get<Goal[]>(GOALS_API_URL);
      setGoals(response.data);
      console.log('Goals fetched:', response.data);
    } catch (error) {
      console.error('Error fetching goals:', error);
      Alert.alert('Error', 'Failed to fetch goals. Make sure your backend is running.');
    } finally {
      setLoading(false);
    }
  };

  // POST: Create a new goal
  const createGoal = async () => {
    if (!title.trim()) {
      Alert.alert('Error', 'Please enter a goal title');
      return;
    }

    try {
      const newGoal = {
        title: title,
        description: description,
        completed: false,
      };

      const response = await axios.post<Goal>(GOALS_API_URL, newGoal);
      setGoals([...goals, response.data]);
      setTitle('');
      setDescription('');
      Alert.alert('Success', 'Goal created!');
    } catch (error) {
      console.error('Error creating goal:', error);
      Alert.alert('Error', 'Failed to create goal');
    }
  };

  // PUT: Update an existing goal
  const updateGoal = async () => {
    if (!editingGoal) return;

    try {
      const updatedGoal = {
        title: title,
        description: description,
        completed: editingGoal.completed,
      };

      const response = await axios.put<Goal>(
        `${GOALS_API_URL}/${editingGoal.id}`,
        updatedGoal
      );

      setGoals(goals.map(goal =>
        goal.id === editingGoal.id ? response.data : goal
      ));

      setTitle('');
      setDescription('');
      setEditingGoal(null);
      Alert.alert('Success', 'Goal updated!');
    } catch (error) {
      console.error('Error updating goal:', error);
      Alert.alert('Error', 'Failed to update goal');
    }
  };

  // DELETE: Delete a goal
  const deleteGoal = async (goalId: string) => {
    try {
      await axios.delete(`${GOALS_API_URL}/${goalId}`);
      setGoals(goals.filter(goal => goal.id !== goalId));
      Alert.alert('Success', 'Goal deleted!');
    } catch (error) {
      console.error('Error deleting goal:', error);
      Alert.alert('Error', 'Failed to delete goal');
    }
  };

  // Toggle goal completion status
  const toggleComplete = async (goal: Goal) => {
    try {
      const updatedGoal = {
        ...goal,
        completed: !goal.completed,
      };

      const response = await axios.put<Goal>(
        `${GOALS_API_URL}/${goal.id}`,
        updatedGoal
      );

      setGoals(goals.map(g =>
        g.id === goal.id ? response.data : g
      ));
    } catch (error) {
      console.error('Error toggling goal:', error);
      Alert.alert('Error', 'Failed to update goal');
    }
  };

  // Start editing a goal
  const startEditing = (goal: Goal) => {
    setEditingGoal(goal);
    setTitle(goal.title);
    setDescription(goal.description);
  };

  // Cancel editing
  const cancelEditing = () => {
    setEditingGoal(null);
    setTitle('');
    setDescription('');
  };

  // Render a single goal item
  const renderGoal = ({ item }: { item: Goal }) => (
    <View style={styles.goalItem}>
      <TouchableOpacity
        style={styles.goalContent}
        onPress={() => toggleComplete(item)}
      >
        <View style={styles.checkbox}>
          {item.completed && <View style={styles.checkboxChecked} />}
        </View>
        <View style={styles.goalText}>
          <Text style={[
            styles.goalTitle,
            item.completed && styles.completedText
          ]}>
            {item.title}
          </Text>
          {item.description ? (
            <Text style={styles.goalDescription}>{item.description}</Text>
          ) : null}
        </View>
      </TouchableOpacity>

      <View style={styles.goalActions}>
        <TouchableOpacity
          style={styles.editButton}
          onPress={() => startEditing(item)}
        >
          <Text style={styles.editButtonText}>Edit</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={styles.deleteButton}
          onPress={() => deleteGoal(item.id)}
        >
          <Text style={styles.deleteButtonText}>Delete</Text>
        </TouchableOpacity>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <Text style={styles.header}>MasteryPath</Text>

      {/* Input Form */}
      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          placeholder="Goal Title"
          value={title}
          onChangeText={setTitle}
        />
        <TextInput
          style={styles.input}
          placeholder="Description (optional)"
          value={description}
          onChangeText={setDescription}
        />

        <View style={styles.buttonRow}>
          {editingGoal ? (
            <>
              <TouchableOpacity
                style={[styles.button, styles.updateButton]}
                onPress={updateGoal}
              >
                <Text style={styles.buttonText}>Update Goal</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={cancelEditing}
              >
                <Text style={styles.buttonText}>Cancel</Text>
              </TouchableOpacity>
            </>
          ) : (
            <TouchableOpacity
              style={styles.button}
              onPress={createGoal}
            >
              <Text style={styles.buttonText}>Add Goal</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Goal List */}
      {loading ? (
        <ActivityIndicator size="large" color="#007AFF" />
      ) : (
        <FlatList
          data={goals}
          renderItem={renderGoal}
          keyExtractor={item => item.id}
          contentContainerStyle={styles.listContainer}
          ListEmptyComponent={
            <Text style={styles.emptyText}>No goals yet. Create one above!</Text>
          }
        />
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    fontSize: 28,
    fontWeight: 'bold',
    textAlign: 'center',
    marginTop: 20,
    marginBottom: 20,
    color: '#333',
  },
  inputContainer: {
    backgroundColor: 'white',
    padding: 15,
    marginHorizontal: 15,
    marginBottom: 20,
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    padding: 12,
    marginBottom: 10,
    borderRadius: 8,
    fontSize: 16,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: 10,
  },
  button: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 8,
    flex: 1,
  },
  updateButton: {
    backgroundColor: '#34C759',
  },
  cancelButton: {
    backgroundColor: '#FF9500',
  },
  buttonText: {
    color: 'white',
    textAlign: 'center',
    fontSize: 16,
    fontWeight: '600',
  },
  listContainer: {
    paddingHorizontal: 15,
    paddingBottom: 20,
  },
  goalItem: {
    backgroundColor: 'white',
    padding: 15,
    marginBottom: 10,
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 2,
  },
  goalContent: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 10,
  },
  checkbox: {
    width: 24,
    height: 24,
    borderRadius: 12,
    borderWidth: 2,
    borderColor: '#007AFF',
    marginRight: 12,
    justifyContent: 'center',
    alignItems: 'center',
  },
  checkboxChecked: {
    width: 14,
    height: 14,
    borderRadius: 7,
    backgroundColor: '#007AFF',
  },
  goalText: {
    flex: 1,
  },
  goalTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  goalDescription: {
    fontSize: 14,
    color: '#666',
    marginTop: 4,
  },
  completedText: {
    textDecorationLine: 'line-through',
    color: '#999',
  },
  goalActions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 10,
  },
  editButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  editButtonText: {
    color: 'white',
    fontSize: 14,
    fontWeight: '600',
  },
  deleteButton: {
    backgroundColor: '#FF3B30',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  deleteButtonText: {
    color: 'white',
    fontSize: 14,
    fontWeight: '600',
  },
  emptyText: {
    textAlign: 'center',
    color: '#999',
    marginTop: 40,
    fontSize: 16,
  },
});