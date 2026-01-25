import React, { useState } from 'react';
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
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

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

// API functions
const fetchGoals = async (): Promise<Goal[]> => {
  const response = await axios.get<Goal[]>(GOALS_API_URL);
  return response.data;
};

const createGoalApi = async (goal: Partial<Goal>): Promise<Goal> => {
  const response = await axios.post<Goal>(GOALS_API_URL, goal);
  return response.data;
};

const updateGoalApi = async (goal: Goal): Promise<Goal> => {
  const response = await axios.put<Goal>(`${GOALS_API_URL}/${goal.id}`, goal);
  return response.data;
};

const deleteGoalApi = async (goalId: string): Promise<void> => {
  await axios.delete(`${GOALS_API_URL}/${goalId}`);
};

export default function MasteryPathScreen() {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null);

  // Query: Fetch all goals
  const { data: goals = [], isLoading, error } = useQuery({
    queryKey: ['goals'],
    queryFn: fetchGoals,
  });

  // Mutation: Create
  const createMutation = useMutation({
    mutationFn: createGoalApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      setTitle('');
      setDescription('');
      Alert.alert('Success', 'Goal created!');
    },
    onError: () => Alert.alert('Error', 'Failed to create goal'),
  });

  // Mutation: Update
  const updateMutation = useMutation({
    mutationFn: updateGoalApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      setTitle('');
      setDescription('');
      setEditingGoal(null);
      Alert.alert('Success', 'Goal updated!');
    },
    onError: () => Alert.alert('Error', 'Failed to update goal'),
  });

  // Mutation: Delete
  const deleteMutation = useMutation({
    mutationFn: deleteGoalApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      Alert.alert('Success', 'Goal deleted!');
    },
    onError: () => Alert.alert('Error', 'Failed to delete goal'),
  });

  // Mutation: Toggle Complete (with Optimistic Update)
  const toggleMutation = useMutation({
    mutationFn: updateGoalApi,
    onMutate: async (updatedGoal) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['goals'] });

      // Snapshot the previous value
      const previousGoals = queryClient.getQueryData<Goal[]>(['goals']);

      // Optimistically update to the new value
      queryClient.setQueryData<Goal[]>(['goals'], (old) =>
        old?.map((g) => (g.id === updatedGoal.id ? updatedGoal : g))
      );

      // Return a context with the snapshotted value
      return { previousGoals };
    },
    onError: (_err, _updatedGoal, context) => {
      // Rollback on error
      queryClient.setQueryData(['goals'], context?.previousGoals);
      Alert.alert('Error', 'Failed to update goal');
    },
    onSettled: () => {
      // Always refetch after error or success
      queryClient.invalidateQueries({ queryKey: ['goals'] });
    },
  });

  const handleCreateGoal = () => {
    if (!title.trim()) {
      Alert.alert('Error', 'Please enter a goal title');
      return;
    }
    createMutation.mutate({
      title,
      description,
      completed: false,
    });
  };

  const handleUpdateGoal = () => {
    if (!editingGoal) return;
    updateMutation.mutate({
      ...editingGoal,
      title,
      description,
    });
  };

  const handleDeleteGoal = (goalId: string) => {
    deleteMutation.mutate(goalId);
  };

  const handleToggleComplete = (goal: Goal) => {
    toggleMutation.mutate({
      ...goal,
      completed: !goal.completed,
    });
  };

  const startEditing = (goal: Goal) => {
    setEditingGoal(goal);
    setTitle(goal.title);
    setDescription(goal.description);
  };

  const cancelEditing = () => {
    setEditingGoal(null);
    setTitle('');
    setDescription('');
  };

  const renderGoal = ({ item }: { item: Goal }) => (
    <View style={styles.goalItem}>
      <TouchableOpacity
        style={styles.goalContent}
        onPress={() => handleToggleComplete(item)}
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
          onPress={() => handleDeleteGoal(item.id)}
        >
          <Text style={styles.deleteButtonText}>Delete</Text>
        </TouchableOpacity>
      </View>
    </View>
  );

  if (error) {
    return (
      <SafeAreaView style={styles.container}>
        <Text style={styles.header}>MasteryPath</Text>
        <Text style={styles.emptyText}>Failed to load goals. Make sure your backend is running.</Text>
      </SafeAreaView>
    );
  }

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
                onPress={handleUpdateGoal}
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
              onPress={handleCreateGoal}
            >
              <Text style={styles.buttonText}>Add Goal</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Goal List */}
      {isLoading ? (
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