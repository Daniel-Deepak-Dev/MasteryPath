import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Colors, getMasteryColor, getMasteryLabel, RadarAxisColors } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';
import RadarChart from './radar-chart';
import type { DashboardSkill } from '@/services/api';

interface SkillCardProps {
  skill: DashboardSkill;
}

export default function SkillCard({ skill }: SkillCardProps) {
  const colorScheme = useColorScheme() ?? 'light';
  const colors = Colors[colorScheme];
  const masteryColor = getMasteryColor(skill.overall_mastery);
  const masteryLabel = getMasteryLabel(skill.overall_mastery);

  const radarData = skill.sub_skills.map((s) => ({
    label: s.name,
    value: s.mastery,
  }));

  return (
    <View style={[styles.card, { backgroundColor: colors.surface, borderColor: colors.border }]}>
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={[styles.skillName, { color: colors.text }]}>{skill.name}</Text>
          {skill.category ? (
            <View style={[styles.categoryBadge, { backgroundColor: colors.chipBg }]}>
              <Text style={[styles.categoryText, { color: colors.chipText }]}>
                {skill.category}
              </Text>
            </View>
          ) : null}
        </View>
        <View style={[styles.masteryBadge, { backgroundColor: masteryColor + '18' }]}>
          <Text style={[styles.masteryPercent, { color: masteryColor }]}>
            {skill.overall_mastery}%
          </Text>
          <Text style={[styles.masteryLabel, { color: masteryColor }]}>{masteryLabel}</Text>
        </View>
      </View>

      {/* Description */}
      {skill.description ? (
        <Text style={[styles.description, { color: colors.textSecondary }]} numberOfLines={2}>
          {skill.description}
        </Text>
      ) : null}

      {/* Radar Chart */}
      {skill.sub_skills.length >= 3 ? (
        <View style={styles.chartContainer}>
          <RadarChart
            data={radarData}
            size={220}
            fillColor={
              colorScheme === 'dark'
                ? 'rgba(167, 139, 250, 0.25)'
                : 'rgba(108, 99, 255, 0.2)'
            }
            strokeColor={colorScheme === 'dark' ? '#A78BFA' : '#6C63FF'}
          />
        </View>
      ) : null}

      {/* Sub-skill chips */}
      <View style={styles.chipRow}>
        {skill.sub_skills.map((sub, index) => {
          const chipColor = RadarAxisColors[index % RadarAxisColors.length];
          return (
            <View
              key={sub.id}
              style={[styles.chip, { backgroundColor: chipColor + '15', borderColor: chipColor + '30' }]}
            >
              <View style={[styles.chipDot, { backgroundColor: chipColor }]} />
              <Text style={[styles.chipLabel, { color: colors.text }]}>{sub.name}</Text>
              <Text style={[styles.chipValue, { color: chipColor }]}>
                {sub.mastery}%
              </Text>
            </View>
          );
        })}
        {skill.sub_skills.length === 0 && (
          <Text style={[styles.emptySubSkills, { color: colors.textSecondary }]}>
            No sub-skills added yet
          </Text>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: 16,
    padding: 20,
    marginHorizontal: 16,
    marginBottom: 16,
    borderWidth: 1,
    // Shadow
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.08,
    shadowRadius: 12,
    elevation: 4,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  headerLeft: {
    flex: 1,
    marginRight: 12,
  },
  skillName: {
    fontSize: 20,
    fontWeight: '700',
    letterSpacing: -0.3,
    marginBottom: 6,
  },
  categoryBadge: {
    alignSelf: 'flex-start',
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 20,
  },
  categoryText: {
    fontSize: 12,
    fontWeight: '600',
  },
  masteryBadge: {
    alignItems: 'center',
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 12,
  },
  masteryPercent: {
    fontSize: 22,
    fontWeight: '800',
    letterSpacing: -0.5,
  },
  masteryLabel: {
    fontSize: 10,
    fontWeight: '600',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginTop: 2,
  },
  description: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 12,
  },
  chartContainer: {
    alignItems: 'center',
    paddingVertical: 8,
  },
  chipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginTop: 4,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 10,
    paddingVertical: 6,
    borderRadius: 20,
    borderWidth: 1,
  },
  chipDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    marginRight: 6,
  },
  chipLabel: {
    fontSize: 12,
    fontWeight: '600',
    marginRight: 4,
  },
  chipValue: {
    fontSize: 12,
    fontWeight: '700',
  },
  emptySubSkills: {
    fontSize: 13,
    fontStyle: 'italic',
    paddingVertical: 4,
  },
});
