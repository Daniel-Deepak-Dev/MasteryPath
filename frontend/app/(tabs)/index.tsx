import React, { useCallback } from 'react';
import {
  StyleSheet,
  Text,
  View,
  FlatList,
  ActivityIndicator,
  RefreshControl,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useQuery } from '@tanstack/react-query';
import { useColorScheme } from '@/hooks/use-color-scheme';
import { Colors } from '@/constants/theme';
import { fetchDashboard, DashboardSkill } from '@/services/api';
import SkillCard from '@/components/skill-card';
import { ContributionGraph } from '@/components/ContributionGraph'; // Add to imports
import axios from 'axios'; // Assuming axios is used for these new queries

// Define API_URL if not already defined globally or in a config file
// For demonstration, using a placeholder. In a real app, this would come from env vars or a config.
const API_URL = 'http://localhost:3000/api'; // Placeholder API URL

export default function DashboardScreen() {
  const colorScheme = useColorScheme() ?? 'light';
  const colors = Colors[colorScheme];

  const {
    data: skills = [],
    isLoading,
    error,
    refetch,
    isRefetching,
  } = useQuery<DashboardSkill[]>({
    queryKey: ['dashboard'],
    queryFn: fetchDashboard,
  });

  // 1. Fetch Goals
  const { data: goals = [], refetch: refetchGoals } = useQuery({
    queryKey: ['goals'],
    queryFn: async () => {
      const response = await axios.get(`${API_URL}/goals`);
      return response.data;
    },
  });

  // 1b. Fetch Contribution Data
  const { data: contributions = [] } = useQuery({
    queryKey: ['contributions'],
    queryFn: async () => {
      try {
         const response = await axios.get(`${API_URL}/analytics/contributions`);
         return response.data || [];
      } catch (e) {
         console.log('Error fetching contributions:', e);
         return [];
      }
    },
  });

  const onRefresh = useCallback(() => {
    refetch();
    refetchGoals(); // Also refresh goals
  }, [refetch, refetchGoals]);

  const renderSkillCard = ({ item }: { item: DashboardSkill }) => (
    <SkillCard skill={item} />
  );

  if (error) {
    return (
      <SafeAreaView style={[styles.container, { backgroundColor: colors.background }]}>
        <View style={styles.headerSection}>
          <Text style={[styles.headerTitle, { color: colors.text }]}>MasteryPath</Text>
          <Text style={[styles.headerSubtitle, { color: colors.textSecondary }]}>
            Your learning journey
          </Text>
        </View>
        <View style={styles.emptyContainer}>
          <Text style={styles.errorEmoji}>⚠️</Text>
          <Text style={[styles.errorTitle, { color: colors.text }]}>Connection Error</Text>
          <Text style={[styles.errorText, { color: colors.textSecondary }]}>
            Could not reach the server.{'\n'}Make sure the backend is running.
          </Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={isLoading}
            onRefresh={onRefresh}
            colors={[colors.tint]}
          />
        }
      >
        {/* Header */}
        <View style={styles.headerSection}>
          <View>
            <Text style={[styles.headerTitle, { color: colors.text }]}>MasteryPath</Text>
            <Text style={[styles.headerSubtitle, { color: colors.textSecondary }]}>
              Your learning journey
            </Text>
          </View>
          {/* Assuming 'router' and 'Ionicons' are imported or globally available */}
          {/* <TouchableOpacity onPress={() => router.push('/modal')} style={styles.addButton}>
            <Ionicons name="add" size={24} color="#FFF" />
          </TouchableOpacity> */}
        </View>

        {/* Contribution Graph */}
        <View style={[styles.section, { backgroundColor: colors.surface, marginTop: 10, marginHorizontal: 20, borderRadius: 12, padding: 15 }]}>
           <Text style={[styles.sectionTitle, { color: colors.text, marginBottom: 10, fontSize: 18, fontWeight: 'bold' }]}>Consistency</Text>
           <ContributionGraph data={contributions} />
        </View>

        {/* Stats Row */}
                  </Text>
                  <Text style={[styles.statLabel, { color: colors.textSecondary }]}>Avg Mastery</Text>
                </View>
              </View>
            )}

            {/* Section title */}
            {skills.length > 0 && (
              <Text style={[styles.sectionTitle, { color: colors.text }]}>
                Skill Domains
              </Text>
            )}
          </View>
        }
        ListEmptyComponent={
          isLoading ? (
            <View style={styles.loadingContainer}>
              <ActivityIndicator size="large" color={colors.tint} />
              <Text style={[styles.loadingText, { color: colors.textSecondary }]}>
                Loading your skills...
              </Text>
            </View>
          ) : (
            <View style={styles.emptyContainer}>
              <Text style={styles.emptyEmoji}>🎯</Text>
              <Text style={[styles.emptyTitle, { color: colors.text }]}>
                No Skill Domains Yet
              </Text>
              <Text style={[styles.emptyText, { color: colors.textSecondary }]}>
                Start by adding your first Master Skill and its sub-skills to begin tracking your
                learning journey.
              </Text>
            </View>
          )
        }
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  listContent: {
    paddingBottom: 24,
  },
  headerSection: {
    paddingHorizontal: 20,
    paddingTop: 16,
    paddingBottom: 8,
  },
  headerTitle: {
    fontSize: 32,
    fontWeight: '800',
    letterSpacing: -0.8,
  },
  headerSubtitle: {
    fontSize: 16,
    marginTop: 4,
    marginBottom: 16,
  },
  statsRow: {
    flexDirection: 'row',
    gap: 10,
    marginBottom: 20,
  },
  statCard: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: 14,
    borderRadius: 12,
    borderWidth: 1,
    // Shadow
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.04,
    shadowRadius: 8,
    elevation: 2,
  },
  statNumber: {
    fontSize: 24,
    fontWeight: '800',
    letterSpacing: -0.5,
  },
  statLabel: {
    fontSize: 12,
    fontWeight: '600',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginTop: 2,
  },
  section: {
    marginBottom: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    marginBottom: 12,
  },
  loadingContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingTop: 80,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 15,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingTop: 60,
    paddingHorizontal: 40,
  },
  emptyEmoji: {
    fontSize: 48,
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: '700',
    marginBottom: 8,
    textAlign: 'center',
  },
  emptyText: {
    fontSize: 15,
    lineHeight: 22,
    textAlign: 'center',
  },
  errorEmoji: {
    fontSize: 48,
    marginBottom: 16,
  },
  errorTitle: {
    fontSize: 20,
    fontWeight: '700',
    marginBottom: 8,
    textAlign: 'center',
  },
  errorText: {
    fontSize: 15,
    lineHeight: 22,
    textAlign: 'center',
  },
});